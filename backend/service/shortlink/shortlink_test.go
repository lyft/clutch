package shortlink

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	shortlinkv1cfg "github.com/lyft/clutch/backend/api/config/service/shortlink/v1"
	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/mock/service/dbmock"
	"github.com/lyft/clutch/backend/service"
)

func TestNewDefaults(t *testing.T) {
	service.Registry["clutch.service.db.postgres"] = dbmock.NewMockDB()
	cfg := &shortlinkv1cfg.Config{}

	anycfg, err := anypb.New(cfg)
	assert.NoError(t, err)

	c, err := New(anycfg, zap.NewNop(), tally.NoopScope)
	assert.NoError(t, err)

	slClient := c.(*client)
	assert.Equal(t, defaultShortlinkChars, slClient.shortlinkChars)
	assert.Equal(t, defaultShortlinkLength, slClient.shortlinkLength)
}

func TestNewWithOverrides(t *testing.T) {
	service.Registry["clutch.service.db.postgres"] = dbmock.NewMockDB()
	cfg := &shortlinkv1cfg.Config{
		ShortlinkChars:  "abc",
		ShortlinkLength: 3,
	}

	anycfg, err := anypb.New(cfg)
	assert.NoError(t, err)

	c, err := New(anycfg, zap.NewNop(), tally.NoopScope)
	assert.NoError(t, err)

	slClient := c.(*client)
	assert.Equal(t, "abc", slClient.shortlinkChars)
	assert.Equal(t, 3, slClient.shortlinkLength)
}

func TestGetShortlink(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	slClient := &client{
		shortlinkChars:  "a",
		shortlinkLength: 1,
		db:              m.DB(),
		log:             zap.NewNop(),
	}

	expectedState := []*shortlinkv1.ShareableState{
		{
			Key: "mock",
			State: &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: "mock string"},
			},
		},
	}

	stateJson, err := marshalShareableState(expectedState)
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"page_path", "state"})
	rows.AddRow("/test", stateJson)

	m.Mock.ExpectQuery("SELECT page_path, state FROM shortlink WHERE slhash = .*").
		WillReturnRows(rows)

	path, actualState, err := slClient.Get(context.TODO(), "test")
	assert.NoError(t, err)
	assert.Equal(t, "/test", path)
	assert.Equal(t, expectedState, actualState)
	m.MustMeetExpectations()
}

func TestCreateShortlinkWithRetries(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	slClient := &client{
		shortlinkChars:  "a",
		shortlinkLength: 1,
		db:              m.DB(),
	}

	m.Mock.ExpectExec("INSERT INTO shortlink").WithArgs(
		"a", "/test", []byte("state"),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	hash, err := slClient.createShortlinkWithRetries(context.TODO(), "/test", []byte("state"))
	assert.NoError(t, err)
	assert.NotNil(t, hash)
	m.MustMeetExpectations()
}

func TestGenerateShortlink(t *testing.T) {
	tests := []struct {
		name         string
		inputChars   string
		inputLength  int
		expectLength int
		shouldError  bool
	}{
		{
			name:         "lower alpha 10 len",
			inputChars:   "abcdefghijklmnopqrstuvwxyz",
			inputLength:  10,
			expectLength: 10,
			shouldError:  false,
		},
		{
			name:         "zero len",
			inputChars:   "",
			inputLength:  10,
			expectLength: 0,
			shouldError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := generateShortlink(test.inputChars, test.inputLength)
			if test.shouldError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.Len(t, hash, test.expectLength)
			}
		})
	}
}

func TestProtoAnyForState(t *testing.T) {
	tests := []struct {
		name   string
		expect string
		input  []*shortlinkv1.ShareableState
	}{
		{
			name:   "string value",
			expect: `{"state":[{"key":"mock","state":"mock string"}]}`,
			input: []*shortlinkv1.ShareableState{
				{
					Key: "mock",
					State: &structpb.Value{
						Kind: &structpb.Value_StringValue{StringValue: "mock string"},
					},
				},
			},
		},
		{
			name:   "numbers",
			expect: `{"state":[{"key":"mock","state":{"key":123,"key1":345}}]}`,
			input: []*shortlinkv1.ShareableState{
				{
					Key: "mock",
					State: &structpb.Value{
						Kind: &structpb.Value_StructValue{
							StructValue: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"key":  structpb.NewNumberValue(123),
									"key1": structpb.NewNumberValue(345),
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "nested object",
			expect: `{"state":[{"key":"mock","state":{"key":true,"key1":"value"}}]}`,
			input: []*shortlinkv1.ShareableState{
				{
					Key: "mock",
					State: &structpb.Value{
						Kind: &structpb.Value_StructValue{
							StructValue: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"key":  structpb.NewBoolValue(true),
									"key1": structpb.NewStringValue("value"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			marshal, err := marshalShareableState(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expect, string(marshal))
		})
	}
}
