package sourcegraph

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/shurcooL/graphql"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	sourcegraphv1cfg "github.com/lyft/clutch/backend/api/config/service/sourcegraph/v1"
	sourcegraphv1 "github.com/lyft/clutch/backend/api/sourcegraph/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.sourcegraph"

type Service interface {
	CompareCommits(context.Context, *sourcegraphv1.CompareCommitsRequest) (*sourcegraphv1.CompareCommitsResponse, error)
}

type QraphQLQuery struct {
	Query string `json:"query"`
}

type client struct {
	config *sourcegraphv1cfg.Config
	log    *zap.Logger
	scope  tally.Scope

	gqlClient *graphql.Client
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

	gqlClient := graphql.NewClient(sgURL.String(), oauthClient)

	return &client{
		config: sgConfig,
		log:    log,
		scope:  scope,

		gqlClient: gqlClient,
	}, nil
}

func (c *client) CompareCommits(ctx context.Context, req *sourcegraphv1.CompareCommitsRequest) (*sourcegraphv1.CompareCommitsResponse, error) {
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

	variables := map[string]interface{}{
		"name": graphql.String(req.Repository),
		"base": graphql.String(req.Base),
		"head": graphql.String(req.Head),
	}

	err := c.gqlClient.Query(context.Background(), &compareCommitsQuery, variables)
	if err != nil {
		c.log.Error("None successful response from sourcegraph", zap.Error(err))
		return nil, errors.New("None successful response from sourcegraph")
	}

	if len(compareCommitsQuery.Repository.Comparison.Commits.Nodes) == 0 {
		return &sourcegraphv1.CompareCommitsResponse{
			Commits: []*sourcegraphv1.Commit{},
		}, nil
	}

	commits := make([]*sourcegraphv1.Commit, len(compareCommitsQuery.Repository.Comparison.Commits.Nodes))
	for i, v := range compareCommitsQuery.Repository.Comparison.Commits.Nodes {
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
