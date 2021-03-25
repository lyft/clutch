package authn

// <!-- START clutchdoc -->
// description: Stores tokens from the auth provider(s) in the database.
// <!-- END clutchdoc -->

import (
	"context"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	"github.com/lyft/clutch/backend/service"
)

const StorageName = "clutch.service.authn.storage"

func NewStorage(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	c := &authnv1.StorageConfig{}
	if err := cfg.UnmarshalTo(c); err != nil {
		return nil, err
	}
	return newStorage(c)
}

type Storage interface {
	Store(ctx context.Context, userID, provider string, token *oauth2.Token) error
	Read(ctx context.Context, userID, provider string) (*oauth2.Token, error)
}

type storage struct {
	crypto *cryptographer
	repo   *repository
}

func newStorage(cfg *authnv1.StorageConfig) (Storage, error) {
	if cfg == nil {
		return nil, nil
	}

	crypto, err := newCryptographer(cfg.EncryptionPassphrase)
	if err != nil {
		return nil, err
	}

	repo, err := newRepository()
	if err != nil {
		return nil, err
	}

	return &storage{
		crypto: crypto,
		repo:   repo,
	}, nil
}

func (s *storage) Store(ctx context.Context, userID, provider string, t *oauth2.Token) error {
	if t == nil {
		return status.Errorf(codes.InvalidArgument, "token provided for storage was nil")
	} else if userID == "" || provider == "" || t.AccessToken == "" {
		return status.Errorf(codes.InvalidArgument, "userID '%s' or provider '%s' were blank", userID, provider)
	}

	at, err := s.crypto.Encrypt([]byte(t.AccessToken))
	if err != nil {
		return err
	}

	dbToken := &authnToken{
		userID:      userID,
		provider:    provider,
		accessToken: at,
		expiry:      t.Expiry,
	}

	// Encrypt and store refresh token if present.
	if t.RefreshToken != "" {
		rt, err := s.crypto.Encrypt([]byte(t.RefreshToken))
		if err != nil {
			return err
		}
		dbToken.refreshToken = rt
	}

	// Encrypt and store ID token if present.
	if it, ok := t.Extra("id_token").(string); ok {
		idToken, err := s.crypto.Encrypt([]byte(it))
		if err != nil {
			return err
		}
		dbToken.idToken = idToken
	}

	err = s.repo.createOrUpdateProviderToken(ctx, dbToken)

	return err
}

func (s *storage) Read(ctx context.Context, userID string, provider string) (*oauth2.Token, error) {
	t, err := s.repo.readProviderToken(ctx, userID, provider)
	if err != nil {
		return nil, err
	}

	at, err := s.crypto.Decrypt(t.accessToken)
	if err != nil {
		return nil, err
	}

	ret := &oauth2.Token{
		AccessToken: string(at),
		Expiry:      t.expiry,
	}

	// Set refresh token if it exists on the db record.
	if len(t.refreshToken) > 0 {
		rt, err := s.crypto.Decrypt(t.refreshToken)
		if err != nil {
			return nil, err
		}
		ret.RefreshToken = string(rt)
	}

	// Set idToken if it exists on the database record.
	if len(t.idToken) > 0 {
		it, err := s.crypto.Decrypt(t.idToken)
		if err != nil {
			return nil, err
		}
		ret = ret.WithExtra(map[string]interface{}{"id_token": string(it)})
	}

	return ret, nil
}
