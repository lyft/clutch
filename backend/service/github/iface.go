package github

import (
	"context"

	githubv3 "github.com/google/go-github/v54/github"
	"github.com/shurcooL/githubv4"
)

// Wrapper around the go-github client because it has a different suggested
// testing strategy (wrap the struct) than most libraries (mock the interface)
// See https://github.com/google/go-github/issues/113#issuecomment-46023864
type v3client struct {
	Issues        v3issues
	Organizations v3organizations
	PullRequests  v3pullrequests
	Repositories  v3repositories
	Search        v3search
	Users         v3users
	Checks        v3checks
}

type v4client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/issues_comments.go.
// Method comments below reproduced directly from original definition linked above.
type v3issues interface {
	// CreateComment creates a new comment on the specified issue.
	CreateComment(
		ctx context.Context,
		owner string,
		repo string,
		number int,
		comment *githubv3.IssueComment,
	) (*githubv3.IssueComment, *githubv3.Response, error)
}

// Interface for structs defined in
// https://github.com/google/go-github/blob/master/github/orgs.go
// &&
// https://github.com/google/go-github/blob/master/github/orgs_members.go
// Method comments below reproduced directly from original definition(s) linked above.
type v3organizations interface {
	// Get fetches an organization by name.
	Get(ctx context.Context, org string) (*githubv3.Organization, *githubv3.Response, error)
	// List the organizations for a user. Passing the empty string will list organizations for the authenticated user.
	List(ctx context.Context, user string, opts *githubv3.ListOptions) ([]*githubv3.Organization, *githubv3.Response, error)
	// GetOrgMembership gets the membership for a user in a specified organization.
	// Passing an empty string for user will get the membership for the
	// authenticated user.
	GetOrgMembership(ctx context.Context, user, org string) (*githubv3.Membership, *githubv3.Response, error)
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/pulls.go.
// ,&&
// https://github.com/google/go-github/blob/master/github/pulls_reviews.go
// Method comments below reproduced directly from original definition linked above.
type v3pullrequests interface {
	// Create a new pull request on the specified repository.
	Create(ctx context.Context, owner string, repo string, pull *githubv3.NewPullRequest) (*githubv3.PullRequest, *githubv3.Response, error)
	// Get a single pull request.
	Get(ctx context.Context, owner string, repo string, number int) (*githubv3.PullRequest, *githubv3.Response, error)
	// ListPullRequestsWithCommit returns pull requests associated with a commit SHA.
	ListPullRequestsWithCommit(ctx context.Context, owner, repo, sha string, opts *githubv3.ListOptions) ([]*githubv3.PullRequest, *githubv3.Response, error)
	// ListReviews lists all reviews on the specified pull request.
	ListReviews(ctx context.Context, owner, repo string, number int, opts *githubv3.ListOptions) ([]*githubv3.PullRequestReview, *githubv3.Response, error)
}

// Interface for structs defined in
// https://github.com/google/go-github/blob/master/github/repos.go
// &&
// https://github.com/google/go-github/blob/master/github/repos_commits.go
// &&
// https://github.com/google/go-github/blob/master/github/repos_contents.go
// Method comments below reproduced directly from original definition(s) linked above.
type v3repositories interface {
	// Create a new repository. If an org is specified, the new repository will be created under that org. If the empty string is specified, it will be created for the authenticated user.
	Create(ctx context.Context, org string, repo *githubv3.Repository) (*githubv3.Repository, *githubv3.Response, error)
	// GetContents can return either the metadata and content of a single file
	// (when path references a file) or the metadata of all the files and/or
	// subdirectories of a directory (when path references a directory). To make it
	// easy to distinguish between both result types and to mimic the API as much
	// as possible, both result types will be returned but only one will contain a
	// value and the other will be nil.
	GetContents(ctx context.Context, owner, repo, path string, opt *githubv3.RepositoryContentGetOptions) (*githubv3.RepositoryContent, []*githubv3.RepositoryContent, *githubv3.Response, error)
	// CompareCommits compares a range of commits with each other.
	CompareCommits(ctx context.Context, owner, repo string, base, head string, opts *githubv3.ListOptions) (*githubv3.CommitsComparison, *githubv3.Response, error)
	// GetCommit fetches the specified commit, including all details about it.
	GetCommit(ctx context.Context, owner, repo, sha string, opts *githubv3.ListOptions) (*githubv3.RepositoryCommit, *githubv3.Response, error)
	// DeleteFile deletes a file from a repository and returns the commit.
	// Requires the blob SHA of the file to be deleted.
	DeleteFile(ctx context.Context, owner, repo, path string, opts *githubv3.RepositoryContentFileOptions) (*githubv3.RepositoryContentResponse, *githubv3.Response, error)
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/search.go.
// Method comments below reproduced directly from original definition linked above.
type v3search interface {
	// Code searches code via various criteria.
	Code(ctx context.Context, query string, opts *githubv3.SearchOptions) (*githubv3.CodeSearchResult, *githubv3.Response, error)
	// Issues searches issues via various criteria.
	Issues(ctx context.Context, query string, opts *githubv3.SearchOptions) (*githubv3.IssuesSearchResult, *githubv3.Response, error)
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/users.go.
// Method comments below reproduced directly from original definition linked above.
type v3users interface {
	// Get fetches a user. Passing the empty string will fetch the authenticated user.
	Get(ctx context.Context, user string) (*githubv3.User, *githubv3.Response, error)
}

// Interface for struct defined in https://github.com/google/go-github/blob/master/github/checks.go
// Method comments below reproduced directly from original definition linked above.
type v3checks interface {
	// Lists check runs for a specific ref.
	ListCheckRunsForRef(ctx context.Context, owner, repo, ref string, opts *githubv3.ListCheckRunsOptions) (*githubv3.ListCheckRunsResults, *githubv3.Response, error)
}
