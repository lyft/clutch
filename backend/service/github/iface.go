package github

import (
	"context"

	githubv3 "github.com/google/go-github/v32/github"
	"github.com/shurcooL/githubv4"
)

// Wrapper around the go-github client because it has a different suggested
// testing strategy (wrap the struct) than most libraries (mock the interface)
// See https://github.com/google/go-github/issues/113#issuecomment-46023864
type v3client struct {
	Repositories v3repositories
	PullRequests v3pullrequests
	Issues       v3issues
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/repos.go.
// Method comments below reproduced directly from original definition linked above.
type v3repositories interface {
	// Create a new repository. If an org is specified, the new repository will be created under that org. If the empty string is specified, it will be created for the authenticated user.
	Create(ctx context.Context, org string, repo *githubv3.Repository) (*githubv3.Repository, *githubv3.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opt *githubv3.RepositoryContentGetOptions) (*githubv3.RepositoryContent, []*githubv3.RepositoryContent, *githubv3.Response, error)
	CompareCommits(ctx context.Context, owner, repo string, base, head string) (*githubv3.CommitsComparison, *githubv3.Response, error)
	GetCommit(ctx context.Context, owner, repo, sha string) (*githubv3.RepositoryCommit, *githubv3.Response, error)
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/pulls.go.
// Method comments below reproduced directly from original definition linked above.
type v3pullrequests interface {
	// Create a new pull request on the specified repository.
	Create(ctx context.Context, owner string, repo string, pull *githubv3.NewPullRequest) (*githubv3.PullRequest, *githubv3.Response, error)
}

type v4client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

type v3issues interface {
	CreateComment(
		ctx context.Context,
		owner string,
		repo string,
		number int,
		comment *githubv3.IssueComment,
	) (*githubv3.IssueComment, *githubv3.Response, error)
}
