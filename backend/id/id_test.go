package id

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	id1 := NewID()
	id2 := NewID()

	assert.True(t, id1.Valid())
	assert.True(t, id2.Valid())

	assert.True(t, id1 < id2,
		"ids should be monotonically increasing. id1: %x, id2: %x", id1, id2)

	assert.True(t, id1.Time().Unix() <= id2.Time().Unix(),
		"id timestamps should be increasing (or the same). id1: %v, id2: %v", id1.Time(), id2.Time())

	assert.Equal(t, id1.Worker(), id2.Worker(),
		"ids should have same worker IDs")

	assert.True(t, id1.Sequence() < id2.Sequence() || id1.Sequence() == math.MaxUint32,
		"id sequences should be monotonically increasing")

	assert.Zero(t, id1.Version(), "id1 version must be 0")
	assert.Zero(t, id2.Version(), "id2 version must be 0")
}

func TestNewIDWithTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	id := NewIDWithTime(now)

	assert.True(t, id.Valid())
	assert.Equal(t, now.Unix(), id.Time().Unix())
}

func TestID_Sequence_Overflow(t *testing.T) {
	// CANNOT BE PARALLEL

	atomic.StoreUint64(&sequence, math.MaxUint64)

	assert.NotPanics(t, func() {
		id := NewID()
		assert.Zero(t, id.Sequence())
	})
}

func TestParseID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in    string
		valid bool
	}{
		{fmt.Sprint(NewID()), true},
		{"not a number", false},
		{"1", false},
		{"128", false},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()

			id, err := ParseID(tc.in)

			if tc.valid {
				assert.NoError(t, err)
				assert.True(t, id.Valid())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestID_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   ID
		ex   bool
	}{
		{"less than 32 bits", ID(math.MaxUint32 - 1), false},
		{"exactly 32 bits", ID(math.MaxUint32), false},
		{"non-zero version bit", ID(math.MaxUint64), false},

		{"min id", ID(math.MaxUint32 + 1), true},
		{"max id", ID(math.MaxUint64 - 1), true},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.ex, tc.in.Valid())
		})
	}
}

func TestID_Validate(t *testing.T) {
	t.Parallel()

	id := NewID()
	assert.NoError(t, id.Validate())

	id = 1
	assert.Error(t, id.Validate())
}

func TestID_Time(t *testing.T) {
	t.Parallel()

	now := time.Now()
	assert.True(t, now.Unix() <= NewID().Time().Unix())
}

func TestID_Worker(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint32(workerID>>workerShift), NewID().Worker())
}

func TestID_Version(t *testing.T) {
	t.Parallel()

	assert.Zero(t, NewID().Version())
}

func TestID_String(t *testing.T) {
	t.Parallel()

	id := NewID()
	str := id.String()

	nid, err := ParseID(str)
	assert.NoError(t, err)
	assert.Equal(t, str, nid.String())
}

func TestID_MarshalJSON(t *testing.T) {
	t.Parallel()

	id := NewID()
	expected := strconv.Quote(id.String())

	b, err := json.Marshal(id)

	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func TestID_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	good := NewID()

	tests := []struct {
		in  string
		ex  ID
		err bool
	}{
		{good.String(), good, false},
		{strconv.Quote(good.String()), good, false},

		{"true", 0, true},
		{"{}", 0, true},
		{strconv.Quote("not a number"), 0, true},
		{"123", 0, true},
		{strconv.Quote("123"), 0, true},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()

			var id ID
			err := json.Unmarshal([]byte(tc.in), &id)

			if tc.err {
				assert.Error(t, err)
				assert.False(t, id.Valid())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ex, id)
				assert.True(t, id.Valid())
			}
		})
	}
}
