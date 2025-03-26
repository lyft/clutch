package sourcegraph

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/shurcooL/graphql"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	sourcegraphv1cfg "github.com/lyft/clutch/backend/api/config/service/sourcegraph/v1"
	sourcegraphv1 "github.com/lyft/clutch/backend/api/sourcegraph/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.sourcegraph"

type Service interface {
	CompareCommits(context.Context, *sourcegraphv1.CompareCommitsRequest) (*sourcegraphv1.CompareCommitsResponse, error)
	GetQueryResultsCount(context.Context, *sourcegraphv1.GetQueryResultsCountRequest) (*sourcegraphv1.GetQueryResultsCountResponse, error)
}

type client struct {
	config *sourcegraphv1cfg.Config
	log    *zap.Logger
	scope  tally.Scope

	gqlClient *graphql.Client
}

var compareCommitsQuery struct {
	Repository struct {
		Comparison struct {
			Commits struct {
				Nodes []struct {
					Message graphql.String
					Oid     graphql.String
					Author  struct {
						Person struct {
							Email       graphql.String
							DisplayName graphql.String
						}
					}
				}
			}
		} `graphql:"comparison(base: $base, head: $head)"`
	} `graphql:"repository(name: $name)"`
}

var getQueryResultsCountQuery struct {
	Search struct {
		Results struct {
			ResultCount graphql.Int
		}
	} `graphql:"search(query: $query)"`
}

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (service.Service, error) {
	sgConfig := &sourcegraphv1cfg.Config{}
	err := cfg.UnmarshalTo(sgConfig)
	if err != nil {
		return nil, err
	}

	sgURL, err := url.Parse(fmt.Sprintf("%s%s", sgConfig.Host, "/.api/graphql"))
	if err != nil {
		log.Error("unable to parse sourcegraph host", zap.Error(err))
		return nil, err
	}

	// setup http client with auth token
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: sgConfig.Token,
		TokenType:   "token",
	}))

	if sgConfig.TimeoutMs != 0 {
		oauthClient.Timeout = time.Duration(sgConfig.TimeoutMs) * time.Millisecond
	}

	gqlClient := graphql.NewClient(sgURL.String(), oauthClient)

	return &client{
		config: sgConfig,
		log:    log,
		scope:  scope,

		gqlClient: gqlClient,
	}, nil
}

func (c *client) CompareCommits(ctx context.Context, req *sourcegraphv1.CompareCommitsRequest) (resp *sourcegraphv1.CompareCommitsResponse, err error) {
	// TODO: Remove this when a fix has been landed to the upstream to fix the panic from index out of bounds
	defer func() {
		if r := recover(); r != nil {
			c.log.Warn("Recovered from panic in gql client query")
			resp = nil
			err = errors.New("unsuccessful response from sourcegraph")
		}
	}()

	variables := map[string]interface{}{
		"name": graphql.String(req.Repository),
		"base": graphql.String(req.Base),
		"head": graphql.String(req.Head),
	}

	query := &compareCommitsQuery
	e := c.gqlClient.Query(ctx, query, variables)
	if e != nil {
		c.log.Error("unsuccessful response from sourcegraph",
			zap.String("repo", req.Repository),
			zap.String("base", req.Base),
			zap.String("head", req.Head),
			zap.Error(e))

		return nil, errors.New("unsuccessful response from sourcegraph")
	}

	nodes := query.Repository.Comparison.Commits.Nodes
	commits := make([]*sourcegraphv1.Commit, len(nodes))
	for i, v := range nodes {
		commits[i] = &sourcegraphv1.Commit{
			Oid:         string(v.Oid),
			Message:     string(v.Message),
			Email:       string(v.Author.Person.Email),
			DisplayName: string(v.Author.Person.DisplayName),
		}
	}

	return &sourcegraphv1.CompareCommitsResponse{
		Commits: commits,
	}, nil
}

func (c *client) GetQueryResultsCount(ctx context.Context, req *sourcegraphv1.GetQueryResultsCountRequest) (resp *sourcegraphv1.GetQueryResultsCountResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			c.log.Warn("Recovered from panic in gql client query")
			resp = nil
			err = errors.New("unsuccessful response from sourcegraph")
		}
	}()

	variables := map[string]interface{}{
		"query": graphql.String(req.Query),
	}

	query := &getQueryResultsCountQuery
	e := c.gqlClient.Query(ctx, query, variables)
	if e != nil {
		c.log.Error("unsuccessful response from sourcegraph",
			zap.String("query", req.Query),
			zap.Error(e))

		return nil, errors.New("unsuccessful response from sourcegraph")
	}

	return &sourcegraphv1.GetQueryResultsCountResponse{
		Count: uint32(query.Search.Results.ResultCount), //nolint
	}, nil
}
