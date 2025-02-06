package githubmock

import (
	"context"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/ptypes/any"
	githubv3 "github.com/google/go-github/v54/github"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

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

func (s svc) DeleteFile(ctx context.Context, ref *github.RemoteRef, path, sha, message string) (*githubv3.RepositoryContentResponse, error) {
	return &githubv3.RepositoryContentResponse{
		Commit: githubv3.Commit{
			SHA: githubv3.String("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"),
		},
	}, nil
}

func (s svc) GetDirectory(ctx context.Context, ref *github.RemoteRef, path string) (*github.Directory, error) {
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

func (s svc) CompareCommits(ctx context.Context, ref *github.RemoteRef, compareSHA string) (*githubv3.CommitsComparison, error) {
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

func (s svc) ListPullRequestsWithCommit(ctx context.Context, ref *github.RemoteRef, sha string, opts *githubv3.ListOptions) ([]*github.PullRequestInfo, error) {
	prNumber := 12345
	return []*github.PullRequestInfo{
		{
			Number:     prNumber,
			HTMLURL:    fmt.Sprintf("https://github.com/%s/%s/pull/%s", ref.RepoOwner, ref.RepoName, strconv.Itoa(prNumber)),
			BranchName: "my-branch",
		},
	}, nil
}

func (s svc) GetPullRequest(ctx context.Context, owner, repo string, number int) (*githubv3.PullRequest, error) {
	return &githubv3.PullRequest{
		Number: githubv3.Int(4242),
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

func (s svc) CreateCommit(ctx context.Context, ref *github.RemoteRef, message string, files github.FileMap) (*github.Commit, error) {
	return &github.Commit{
		SHA: "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
	}, nil
}

func (s *svc) SearchCode(ctx context.Context, query string, opts *githubv3.SearchOptions) (*githubv3.CodeSearchResult, error) {
	var codeResults []*githubv3.CodeResult

	codeResults = append(codeResults, &githubv3.CodeResult{
		Name:       githubv3.String("file.go"),
		Path:       githubv3.String("path/to/file.go"),
		SHA:        githubv3.String("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"),
		HTMLURL:    githubv3.String("https://github.com/owner/repo/path/to/file.go"),
		Repository: nil,
	})
	return &githubv3.CodeSearchResult{
		Total:             githubv3.Int(1),
		IncompleteResults: githubv3.Bool(false),
		CodeResults:       codeResults,
	}, nil
}

func (s *svc) GetFileContents(ctx context.Context, ref *github.RemoteRef, path string) (*githubv3.RepositoryContent, error) {
	return &githubv3.RepositoryContent{
		Name:     githubv3.String("README.md"),
		Path:     githubv3.String("README.md"),
		Content:  githubv3.String("# Hello World"),
		Encoding: githubv3.String(""),
	}, nil
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}
