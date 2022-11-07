package authn

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	"github.com/lyft/clutch/backend/mock/service/dbmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

type notEmptyBytes struct{}

func (a notEmptyBytes) Match(v driver.Value) bool {
	b, ok := v.([]byte)
	if !ok || len(b) == 0 {
		return false
	}
	return true
}

type emptyBytes struct{}

func (e emptyBytes) Match(v driver.Value) bool {
	b, ok := v.([]byte)
	return ok && len(b) == 0
}

func TestNewStorage(t *testing.T) {
	{
		s, err := newStorage(nil)
		assert.NoError(t, err)
		assert.Nil(t, s)
	}

	{
		dbmock.NewMockDB().Register()
		s, err := newStorage(&authnv1.StorageConfig{
			EncryptionPassphrase: "test",
		})
		assert.NoError(t, err)
		assert.NotNil(t, s)

		assert.NotNil(t, s.(*storage).crypto)
		assert.NotNil(t, s.(*storage).repo)
	}
}

func TestStoreErrors(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Null Token
	err = s.Store(context.Background(), "user@example.com", "clutch.example.com", nil)
	assert.Error(t, err)

	// Empty user
	err = s.Store(context.Background(), "", "clutch.example.com", &oauth2.Token{})
	assert.Error(t, err)

	// Empty provider
	err = s.Store(context.Background(), "", "", &oauth2.Token{})
	assert.Error(t, err)
}

func TestStoreNoIDToken(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	now := time.Now().UTC()
	tok := &oauth2.Token{
		AccessToken:  "a",
		RefreshToken: "r",
		Expiry:       now,
	}

	m.Mock.ExpectExec("INSERT INTO authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com", notEmptyBytes{}, notEmptyBytes{}, ([]byte)(nil), now,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.Store(context.Background(), "user@example.com", "clutch.example.com", tok)
	assert.NoError(t, err)

	m.MustMeetExpectations()
}

func TestStoreWithIDToken(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	now := time.Now().UTC()
	tok := &oauth2.Token{
		AccessToken:  "a",
		RefreshToken: "r",
		Expiry:       now,
	}
	tok = tok.WithExtra(map[string]interface{}{"id_token": "i"})

	m.Mock.ExpectExec("INSERT INTO authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com", notEmptyBytes{}, notEmptyBytes{}, notEmptyBytes{}, now,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.Store(context.Background(), "user@example.com", "clutch.example.com", tok)
	assert.NoError(t, err)

	m.MustMeetExpectations()
}

func TestStoreWithoutRefreshToken(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	now := time.Now().UTC()
	tok := &oauth2.Token{
		AccessToken: "a",
		Expiry:      now,
	}
	tok = tok.WithExtra(map[string]interface{}{"id_token": "i"})

	m.Mock.ExpectExec("INSERT INTO authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com", notEmptyBytes{}, emptyBytes{}, notEmptyBytes{}, now,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.Store(context.Background(), "user@example.com", "clutch.example.com", tok)
	assert.NoError(t, err)

	m.MustMeetExpectations()
}

func TestReadNoRows(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	rows := sqlmock.NewRows([]string{"user_id", "provider", "access_token", "refresh_token", "id_token", "expiry"})

	m.Mock.ExpectQuery("SELECT .*? FROM authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com",
		).WillReturnRows(rows)

	tok, err := s.Read(context.Background(), "user@example.com", "clutch.example.com")
	assert.Error(t, err)
	assert.Nil(t, tok)

	m.MustMeetExpectations()
}

func TestReadWithResult(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	now := time.Now().UTC()

	at, _ := s.(*storage).crypto.Encrypt([]byte("Access"))
	rt, _ := s.(*storage).crypto.Encrypt([]byte("REFRESH"))
	it, _ := s.(*storage).crypto.Encrypt([]byte("id"))

	rows := sqlmock.NewRows([]string{"user_id", "provider", "access_token", "refresh_token", "id_token", "expiry"})
	rows.AddRow("user@example.com", "clutch.example.com", at, rt, it, now)

	m.Mock.ExpectQuery("SELECT .*? FROM authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com",
		).WillReturnRows(rows)

	tok, err := s.Read(context.Background(), "user@example.com", "clutch.example.com")
	assert.NoError(t, err)
	assert.NotNil(t, tok)

	assert.Equal(t, tok.AccessToken, "Access")
	assert.Equal(t, tok.RefreshToken, "REFRESH")
	assert.Equal(t, tok.Extra("id_token"), "id")
	assert.Equal(t, tok.Expiry, now)

	m.MustMeetExpectations()
}

func TestReadWithoutRefreshOrIDToken(t *testing.T) {
	m := dbmock.NewMockDB()
	m.Register()

	s, err := newStorage(&authnv1.StorageConfig{EncryptionPassphrase: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, s)

	now := time.Now().UTC()

	at, _ := s.(*storage).crypto.Encrypt([]byte("Access"))

	rows := sqlmock.NewRows([]string{"user_id", "provider", "access_token", "refresh_token", "id_token", "expiry"})
	rows.AddRow("user@example.com", "clutch.example.com", at, nil, nil, now)

	m.Mock.ExpectQuery("SELECT .*? FROM authn_tokens").
		WithArgs(
			"user@example.com", "clutch.example.com",
		).WillReturnRows(rows)

	tok, err := s.Read(context.Background(), "user@example.com", "clutch.example.com")
	assert.NoError(t, err)
	assert.NotNil(t, tok)

	assert.Equal(t, tok.AccessToken, "Access")
	assert.Empty(t, tok.RefreshToken)
	assert.Empty(t, tok.Extra("id_token"))
	assert.Equal(t, tok.Expiry, now)

	m.MustMeetExpectations()
}
