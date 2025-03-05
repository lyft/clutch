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

func (s *svc) ListCheckRunsForRef(ctx context.Context, ref *github.RemoteRef, opts *githubv3.ListCheckRunsOptions) (*githubv3.ListCheckRunsResults, error) {
	var checkRuns []*githubv3.CheckRun

	checkRuns = append(checkRuns, &githubv3.CheckRun{
		Status:     githubv3.String("completed"),
		Conclusion: githubv3.String("success"),
	})

	return &githubv3.ListCheckRunsResults{
		Total:     githubv3.Int(1),
		CheckRuns: checkRuns,
	}, nil
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) SearchIssues(ctx context.Context, query string, opts *githubv3.SearchOptions) (*githubv3.IssuesSearchResult, error) {
	var issuesResults []*githubv3.Issue

	issuesResults = append(issuesResults, &githubv3.Issue{
		URL:           githubv3.String("https://api.github.com/repos/batterseapower/pinyin-toolkit/issues/132"),
		RepositoryURL: githubv3.String("https://api.github.com/repos/batterseapower/pinyin-toolkit"),
		LabelsURL:     githubv3.String("https://api.github.com/repos/batterseapower/pinyin-toolkit/issues/132/labels{/name}"),
		CommentsURL:   githubv3.String("https://api.github.com/repos/batterseapower/pinyin-toolkit/issues/132/comments"),
		EventsURL:     githubv3.String("https://api.github.com/repos/batterseapower/pinyin-toolkit/issues/132/events"),
		HTMLURL:       githubv3.String("https://github.com/batterseapower/pinyin-toolkit/issues/132"),
		ID:            githubv3.Int64(35802),
		NodeID:        githubv3.String("MDU6SXNzdWUzNTgwMg=="),
		Number:        githubv3.Int(132),
		Title:         githubv3.String("Line Number Indexes Beyond 20 Not Displayed"),
		User: &githubv3.User{
			Login:             githubv3.String("Nick3C"),
			ID:                githubv3.Int64(90254),
			NodeID:            githubv3.String("MDQ6VXNlcjkwMjU0"),
			AvatarURL:         githubv3.String("https://secure.gravatar.com/avatar/934442aadfe3b2f4630510de416c5718?d=https://a248.e.akamai.net/assets.github.com%2Fimages%2Fgravatars%2Fgravatar-user-420.png"),
			GravatarID:        githubv3.String(""),
			URL:               githubv3.String("https://api.github.com/users/Nick3C"),
			HTMLURL:           githubv3.String("https://github.com/Nick3C"),
			FollowersURL:      githubv3.String("https://api.github.com/users/Nick3C/followers"),
			FollowingURL:      githubv3.String("https://api.github.com/users/Nick3C/following{/other_user}"),
			GistsURL:          githubv3.String("https://api.github.com/users/Nick3C/gists{/gist_id}"),
			StarredURL:        githubv3.String("https://api.github.com/users/Nick3C/starred{/owner}{/repo}"),
			SubscriptionsURL:  githubv3.String("https://api.github.com/users/Nick3C/subscriptions"),
			OrganizationsURL:  githubv3.String("https://api.github.com/users/Nick3C/orgs"),
			ReposURL:          githubv3.String("https://api.github.com/users/Nick3C/repos"),
			EventsURL:         githubv3.String("https://api.github.com/users/Nick3C/events{/privacy}"),
			ReceivedEventsURL: githubv3.String("https://api.github.com/users/Nick3C/received_events"),
			Type:              githubv3.String("User"),
			SiteAdmin:         githubv3.Bool(true),
		},
		Labels:   nil,
		State:    githubv3.String("open"),
		Assignee: nil,
		Milestone: &githubv3.Milestone{
			URL:         githubv3.String("https://api.github.com/repos/octocat/Hello-World/milestones/1"),
			HTMLURL:     githubv3.String("https://github.com/octocat/Hello-World/milestones/v1.0"),
			LabelsURL:   githubv3.String("https://api.github.com/repos/octocat/Hello-World/milestones/1/labels"),
			ID:          githubv3.Int64(1002604),
			NodeID:      githubv3.String("MDk6TWlsZXN0b25lMTAwMjYwNA=="),
			Number:      githubv3.Int(1),
			State:       githubv3.String("open"),
			Title:       githubv3.String("v1.0"),
			Description: githubv3.String("Tracking milestone for version 1.0"),
			Creator: &githubv3.User{
				Login:             githubv3.String("octocat"),
				ID:                githubv3.Int64(1),
				NodeID:            githubv3.String("MDQ6VXNlcjE="),
				AvatarURL:         githubv3.String("https://github.com/images/error/octocat_happy.gif"),
				GravatarID:        githubv3.String(""),
				URL:               githubv3.String("https://api.github.com/users/octocat"),
				HTMLURL:           githubv3.String("https://github.com/octocat"),
				FollowersURL:      githubv3.String("https://api.github.com/users/octocat/followers"),
				FollowingURL:      githubv3.String("https://api.github.com/users/octocat/following{/other_user}"),
				GistsURL:          githubv3.String("https://api.github.com/users/octocat/gists{/gist_id}"),
				StarredURL:        githubv3.String("https://api.github.com/users/octocat/starred{/owner}{/repo}"),
				SubscriptionsURL:  githubv3.String("https://api.github.com/users/octocat/subscriptions"),
				OrganizationsURL:  githubv3.String("https://api.github.com/users/octocat/orgs"),
				ReposURL:          githubv3.String("https://api.github.com/users/octocat/repos"),
				EventsURL:         githubv3.String("https://api.github.com/users/octocat/events{/privacy}"),
				ReceivedEventsURL: githubv3.String("https://api.github.com/users/octocat/received_events"),
				Type:              githubv3.String("User"),
				SiteAdmin:         githubv3.Bool(false),
			},
			OpenIssues:   githubv3.Int(4),
			ClosedIssues: githubv3.Int(8),
			CreatedAt:    nil,
			UpdatedAt:    nil,
			ClosedAt:     nil,
			DueOn:        nil,
		},
		Comments:  githubv3.Int(15),
		CreatedAt: nil,
		UpdatedAt: nil,
		ClosedAt:  nil,
		PullRequestLinks: &githubv3.PullRequestLinks{
			URL:      githubv3.String("https://api/github.com/repos/octocat/Hello-World/pull/1347"),
			HTMLURL:  githubv3.String("https://github.com/octocat/Hello-World/pull/1347"),
			DiffURL:  githubv3.String("https://github.com/octocat/Hello-World/pull/1347.diff"),
			PatchURL: githubv3.String("https://api.github.com/repos/octocat/Hello-World/pulls/1347"),
		},
		Body:              githubv3.String("..."),
		Locked:            githubv3.Bool(true),
		AuthorAssociation: githubv3.String("COLLABORATOR"),
		StateReason:       githubv3.String("completed"),
	})

	return &githubv3.IssuesSearchResult{
		Total:             githubv3.Int(280),
		IncompleteResults: githubv3.Bool(false),
		Issues:            issuesResults,
	}, nil
}
