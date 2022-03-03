package shortlink

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/lib/pq"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	shortlinkv1cfg "github.com/lyft/clutch/backend/api/config/service/shortlink/v1"
	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const (
	Name                  = "clutch.service.shortlink"
	defaultHashChars      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	defaultHashLength     = 10
	maxHashCollisionRetry = 5

	// If we hit a key collision inserting a duplicate hash this error is thrown.
	// We catch this error and retry up to the maxHashCollisionRetry limit.
	// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
	pgUniqueErrorCode = "23505"
)

type Service interface {
	Create(context.Context, string, []*shortlinkv1.ShareableState) (string, error)
	Get(context.Context, string) (string, []*shortlinkv1.ShareableState, error)
}

type client struct {
	hashChars  string
	hashLength int

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
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

	hashChars := defaultHashChars
	if slConfig.HashChars != "" {
		hashChars = slConfig.HashChars
	}

	hashLength := defaultHashLength
	if slConfig.HashLength > 0 {
		hashLength = int(slConfig.HashLength)
	}

	c := &client{
		hashChars:  hashChars,
		hashLength: hashLength,
		db:         dbClient.DB(),
		log:        logger,
		scope:      scope,
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
	for i := 0; i < maxHashCollisionRetry; i++ {
		hash, err := generateShortlinkHash(c.hashChars, c.hashLength)
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
		Select("page_path, state").
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

// generateShortlinkHash generates a hash from a set of characters to the length specified
func generateShortlinkHash(chars string, length int) (string, error) {
	if len(chars) == 0 || length == 0 {
		return "", errors.New("chars or length are invalid lengths")
	}

	hash := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		hash[i] = chars[num.Int64()]
	}

	return string(hash), nil
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
