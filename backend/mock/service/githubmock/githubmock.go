package githubmock

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	githubv3 "github.com/google/go-github/v33/github"
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

func (s svc) CreatePullRequest(ctx context.Context, ref *github.RemoteRef, base, title, body string) (*github.PullRequestInfo, error) {
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

func (s svc) GetRepository(ctx context.Context, ref *github.RemoteRef) (*github.Repository, error) {
	return &github.Repository{}, nil
}

func (s svc) GetOrganization(ctx context.Context, organization string) (*githubv3.Organization, error) {
	login := "clutch"
	return &githubv3.Organization{Login: &login}, nil
}

func (s svc) ListOrganizations(ctx context.Context, user string) ([]*githubv3.Organization, error) {
	login := "clutch"
	return []*githubv3.Organization{
		{Login: &login},
	}, nil
}

func (s svc) GetOrgMembership(ctx context.Context, user, org string) (*githubv3.Membership, error) {
	role := "member"
	return &githubv3.Membership{Role: &role}, nil
}

func (s svc) GetUser(ctx context.Context, username string) (*githubv3.User, error) {
	login := "user"
	avatarURL := "https://clutch.sh/img/microsite/logo.svg"
	return &githubv3.User{Login: &login, AvatarURL: &avatarURL}, nil
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
