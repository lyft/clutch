package sourcegraph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	sourcegraphv1cfg "github.com/lyft/clutch/backend/api/config/service/sourcegraph/v1"
	sourcegraphv1 "github.com/lyft/clutch/backend/api/sourcegraph/v1"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestNew(t *testing.T) {
	cfg, _ := anypb.New(&sourcegraphv1cfg.Config{
		Host:  "https://localhost",
		Token: "secret",
	})

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	_, err := New(cfg, log, scope)
	assert.NoError(t, err)
}

func TestCompareCommits(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	tests := []struct {
		id      string
		handler func(http.ResponseWriter, *http.Request)
		req     *sourcegraphv1.CompareCommitsRequest
		res     *sourcegraphv1.CompareCommitsResponse
	}{
		{
			id:  "single response",
			req: &sourcegraphv1.CompareCommitsRequest{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(compareCommitsSingleResponse))
			},
			res: &sourcegraphv1.CompareCommitsResponse{
				Commits: []*sourcegraphv1.Commit{
					{
						Oid:         "8a9857493108d0be9ebc251f30f72665c79424f6",
						Email:       "29139614+renovate[bot]@users.noreply.github.com",
						Message:     "housekeeping: Update dependency cypress to v8.2.0 (#1672)\n",
						DisplayName: "renovate[bot]",
					},
				},
			},
		},
		{
			id:  "multi result response",
			req: &sourcegraphv1.CompareCommitsRequest{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(compareCommitsMultiResultResponse))
			},
			res: &sourcegraphv1.CompareCommitsResponse{
				Commits: []*sourcegraphv1.Commit{
					{
						Oid:         "8a9857493108d0be9ebc251f30f72665c79424f6",
						Email:       "29139614+renovate[bot]@users.noreply.github.com",
						Message:     "housekeeping: Update dependency cypress to v8.2.0 (#1672)\n",
						DisplayName: "renovate[bot]",
					},
					{
						Oid:         "5802372fe18c37e9ab62341b18d7f286078430b4",
						Email:       "29139614+renovate[bot]@users.noreply.github.com",
						Message:     "housekeeping: Update dependency remark-toc to v8 (#1680)\n",
						DisplayName: "renovate[bot]",
					},
				},
			},
		},
		{
			id:  "no results response",
			req: &sourcegraphv1.CompareCommitsRequest{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(compareCommitsNoResultsResponse))
			},
			res: &sourcegraphv1.CompareCommitsResponse{
				Commits: []*sourcegraphv1.Commit{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer srv.Close()

			sgURL, err := url.Parse(srv.URL)
			assert.NoError(t, err)

			c := &client{
				log:   log,
				scope: scope,

				client: srv.Client(),
				sgURL:  sgURL,
			}

			res, err := c.CompareCommits(context.Background(), tt.req)
			assert.NoError(t, err)
			assert.Equal(t, tt.res.Commits, res.Commits)
		})
	}
}

var compareCommitsSingleResponse string = `{"data":{"repository":{"comparison":{"commits":{"nodes":[{"message":"housekeeping: Update dependency cypress to v8.2.0 (#1672)\n","oid":"8a9857493108d0be9ebc251f30f72665c79424f6","author":{"person":{"name":"renovate[bot]","email":"29139614+renovate[bot]@users.noreply.github.com","displayName":"renovate[bot]","avatarURL":null}}}]}}}}}`
var compareCommitsMultiResultResponse string = `{"data":{"repository":{"comparison":{"commits":{"nodes":[{"message":"housekeeping: Update dependency cypress to v8.2.0 (#1672)\n","oid":"8a9857493108d0be9ebc251f30f72665c79424f6","author":{"person":{"name":"renovate[bot]","email":"29139614+renovate[bot]@users.noreply.github.com","displayName":"renovate[bot]","avatarURL":null}}},{"message":"housekeeping: Update dependency remark-toc to v8 (#1680)\n","oid":"5802372fe18c37e9ab62341b18d7f286078430b4","author":{"person":{"name":"renovate[bot]","email":"29139614+renovate[bot]@users.noreply.github.com","displayName":"renovate[bot]","avatarURL":null}}}]}}}}}`
var compareCommitsNoResultsResponse string = `{"data":{"repository":{"comparison":{"commits":{"nodes":[]}}}}}`
