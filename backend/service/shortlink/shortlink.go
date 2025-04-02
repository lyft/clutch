package shortlink

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	shortlinkv1cfg "github.com/lyft/clutch/backend/api/config/service/shortlink/v1"
	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const (
	Name                   = "clutch.service.shortlink"
	defaultShortlinkChars  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	defaultShortlinkLength = 10
	maxCollisionRetry      = 5

	// If we hit a key collision inserting a duplicate random string this error is thrown.
	// We catch this error and retry up to the maxCollisionRetry limit.
	// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
	pgUniqueErrorCode = "23505"
)

type Service interface {
	Create(context.Context, string, []*shortlinkv1.ShareableState) (string, error)
	Get(context.Context, string) (string, []*shortlinkv1.ShareableState, error)
}

type client struct {
	shortlinkChars  string
	shortlinkLength int

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	slConfig := &shortlinkv1cfg.Config{}
	err := cfg.UnmarshalTo(slConfig)
	if err != nil {
		return nil, err
	}

	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("Please config the datastore [clutch.service.db.postgres] to use the shortlink service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get the datastore client")
	}

	chars := defaultShortlinkChars
	if slConfig.ShortlinkChars != "" {
		chars = slConfig.ShortlinkChars
	}

	length := defaultShortlinkLength
	if slConfig.ShortlinkLength > 0 {
		length = int(slConfig.ShortlinkLength)
	}

	c := &client{
		shortlinkChars:  chars,
		shortlinkLength: length,
		db:              dbClient.DB(),
		log:             logger,
		scope:           scope,
	}

	return c, nil
}

func (c *client) Create(ctx context.Context, path string, state []*shortlinkv1.ShareableState) (string, error) {
	stateJson, err := marshalShareableState(state)
	if err != nil {
		return "", err
	}

	return c.createShortlinkWithRetries(ctx, path, stateJson)
}

// createShortlinkWithRetries retries the insert of a new shortlink
// This function generates a shortlink hash which is used as the primary key in the shortlink table
// There could be a possibility of a collision depending on the configuration
// With the default settings in place [a-zA-Z0-9] and a default subset length of 10,
// this leaves us with 62^10.
func (c *client) createShortlinkWithRetries(ctx context.Context, path string, state []byte) (string, error) {
	for i := 0; i < maxCollisionRetry; i++ {
		hash, err := generateShortlink(c.shortlinkChars, c.shortlinkLength)
		if err != nil {
			return "", err
		}

		insertBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Insert("shortlink").
			Columns("slhash", "page_path", "state").
			Values(hash, path, state)

		_, err = insertBuilder.RunWith(c.db).ExecContext(ctx)
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgUniqueErrorCode {
				// If we hit a key collision lets retry
				continue
			} else {
				return "", err
			}
		}

		return hash, err
	}

	return "", errors.New("retries exhausted, unable to create unique shortlink hash.")
}

func (c *client) Get(ctx context.Context, hash string) (string, []*shortlinkv1.ShareableState, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("page_path", "state").
		From("shortlink").
		Where(sq.Eq{"slhash": hash})

	row := query.RunWith(c.db).QueryRowContext(ctx)

	var path string
	var byteState []byte
	if err := row.Scan(&path, &byteState); err != nil {
		c.log.Error("Error scanning row", zap.Error(err))
		return "", nil, err
	}

	var state shortlinkv1.CreateRequest
	if err := protojson.Unmarshal(byteState, &state); err != nil {
		c.log.Error("Error unmarshaling data field", zap.Error(err))
		return "", nil, err
	}

	return path, state.State, nil
}

// generateShortlink generates a random string from a set of characters to the length specified
func generateShortlink(chars string, length int) (string, error) {
	if len(chars) == 0 || length == 0 {
		return "", errors.New("chars or length are invalid lengths")
	}

	res := make([]byte, length)
	_, err := rand.Read(res)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		res[i] = chars[int(res[i])%len(chars)]
	}

	return string(res), nil
}

func marshalShareableState(state []*shortlinkv1.ShareableState) ([]byte, error) {
	stateJson, err := protojson.Marshal(&shortlinkv1.CreateRequest{
		State: state,
	})
	if err != nil {
		return nil, err
	}
	return stateJson, nil
}
