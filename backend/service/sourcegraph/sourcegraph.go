package sourcegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/ptypes/any"
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

	sgURL  *url.URL
	client *http.Client
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
	sgClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: sgConfig.Token,
		TokenType:   "token",
	}))

	return &client{
		config: sgConfig,
		log:    log,
		scope:  scope,

		sgURL:  sgURL,
		client: sgClient,
	}, nil
}

func (c *client) CompareCommits(ctx context.Context, req *sourcegraphv1.CompareCommitsRequest) (*sourcegraphv1.CompareCommitsResponse, error) {
	compareCommitQuery := fmt.Sprintf(`
	query {
		repository(name: "%s") {
			comparison(
				base: "%s"
				head: "%s"
			) {
				commits {
					nodes {
						message
						oid
						author {
							person {
								email
								displayName
							}
						}
					}
				}
			}
		}
	}
	`, req.Repository, req.Base, req.Head)

	gqlQuery := &QraphQLQuery{
		Query: compareCommitQuery,
	}

	data, err := json.Marshal(gqlQuery)
	if err != nil {
		c.log.Error("unable to marshal gql query", zap.Error(err))
		return nil, err
	}

	sgRequest := &http.Request{
		Method: http.MethodPost,
		URL:    c.sgURL,
		Body:   ioutil.NopCloser(bytes.NewBuffer(data)),
	}

	res, err := c.client.Do(sgRequest)
	if err != nil {
		c.log.Error("error querying sourcegraph", zap.Error(err))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		bodyData, err := io.ReadAll(res.Body)
		if err == nil {
			c.log.Error("None successful response from sourcegraph",
				zap.Int("statusCode", res.StatusCode), zap.String("error", string(bodyData)))
		}

		c.log.Error("None successful response from sourcegraph and unable to parse the body",
			zap.Int("statusCode", res.StatusCode))
		return nil, errors.New("None successful response from sourcegraph")
	}

	var bodyData compareCommitsResponse
	err = json.NewDecoder(res.Body).Decode(&bodyData)
	if err != nil {
		c.log.Error("Unable to decode sourcegraph response", zap.Error(err))
		return nil, err
	}

	if len(bodyData.Data.Repository.Comparision.Commits.Nodes) == 0 {
		return &sourcegraphv1.CompareCommitsResponse{
			Commits: []*sourcegraphv1.Commit{},
		}, nil
	}

	commits := make([]*sourcegraphv1.Commit, len(bodyData.Data.Repository.Comparision.Commits.Nodes))
	for i, v := range bodyData.Data.Repository.Comparision.Commits.Nodes {
		commits[i] = &sourcegraphv1.Commit{
			Oid:         v.Oid,
			Message:     v.Message,
			Email:       v.Author.Person.Email,
			DisplayName: v.Author.Person.DisplayName,
		}
	}

	return &sourcegraphv1.CompareCommitsResponse{
		Commits: commits,
	}, nil
}
