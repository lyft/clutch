package githubmock

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	githubv1 "github.com/lyft/clutch/backend/api/sourcecontrol/github/v1"
	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/github"
)

func New() github.Client {
	return &svc{}
}

type svc struct{}

func (s svc) GetFile(ctx context.Context, ref *github.RemoteRef, path string) (*github.File, error) {
	panic("implement me")
}

func (s svc) CreateBranch(ctx context.Context, req *github.CreateBranchRequest) error {
	panic("implement me")
}

func (s svc) CreatePullRequest(ctx context.Context, ref *github.RemoteRef, title, body string) (*github.PullRequestInfo, error) {
	panic("implement me")
}

func (s svc) CreateIssueComment(ctx context.Context, ref *github.RemoteRef, number int, body string) error {
	panic("implement me")
}

func (s svc) CreateRepository(ctx context.Context, req *sourcecontrolv1.CreateRepositoryRequest) (*sourcecontrolv1.CreateRepositoryResponse, error) {
	return &sourcecontrolv1.CreateRepositoryResponse{Url: "https://github.com/lyft/clutch"}, nil
}

func (s svc) CompareCommits(ctx context.Context, ref *github.RemoteRef, compareSHA string) (*githubv1.CommitComparison, error) {
	panic("implement me")
}

func (s svc) GetCommit(ctx context.Context, ref *github.RemoteRef) (*github.Commit, error) {
	panic("implement me")
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
