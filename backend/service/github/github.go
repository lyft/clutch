package github

// <!-- START clutchdoc -->
// description: GitHub client that combines the REST/GraphQL APIs and raw git capabilities into a single interface.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gittransport "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	githubv3 "github.com/google/go-github/v32/github"
	"github.com/shurcooL/githubv4"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	githubv1 "github.com/lyft/clutch/backend/api/config/service/github/v1"
	scgithubv1 "github.com/lyft/clutch/backend/api/sourcecontrol/github/v1"
	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.github"

type FileMap map[string]io.ReadCloser

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &githubv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}
	return newService(config), nil
}

// Remote ref points to a git reference using a combination of the repository and the reference itself.
type RemoteRef struct {
	// Organization or user that owns the repository.
	RepoOwner string
	// Name of the repository.
	RepoName string
	// SHA, branch name, or tag.
	Ref string
}

// File contains information about a requested file, including its content.
type File struct {
	Path             string
	Contents         io.ReadCloser
	SHA              string
	LastModifiedTime time.Time
	LastModifiedSHA  string
}

// Client allows various interactions with remote repositories on GitHub.
type Client interface {
	GetFile(ctx context.Context, ref *RemoteRef, path string) (*File, error)
	CreateBranch(ctx context.Context, req *CreateBranchRequest) error
	CreatePullRequest(ctx context.Context, ref *RemoteRef, title, body string) (*PullRequestInfo, error)
	CreateRepository(ctx context.Context, req *sourcecontrolv1.CreateRepositoryRequest) (*sourcecontrolv1.CreateRepositoryResponse, error)
	CreateIssueComment(ctx context.Context, ref *RemoteRef, number int, body string) error
	CompareCommits(ctx context.Context, ref *RemoteRef, compareSHA string) (*scgithubv1.CommitComparison, error)
	GetCommit(ctx context.Context, ref *RemoteRef) (*Commit, error)
}

// This func can be used to create comments for PRs or Issues
func (s *svc) CreateIssueComment(ctx context.Context, ref *RemoteRef, number int, body string) error {
	com := &githubv3.IssueComment{
		Body: strPtr(body),
	}
	_, _, err := s.rest.Issues.CreateComment(ctx, ref.RepoOwner, ref.RepoName, number, com)
	return err
}

type PullRequestInfo struct {
	Number  int
	HTMLURL string
}

type svc struct {
	graphQL v4client
	rest    v3client
	rawAuth *gittransport.BasicAuth
}

func (s *svc) CreateRepository(ctx context.Context, req *sourcecontrolv1.CreateRepositoryRequest) (*sourcecontrolv1.CreateRepositoryResponse, error) {
	// Validate that we received GitHub Options.
	_, ok := req.Options.(*sourcecontrolv1.CreateRepositoryRequest_GithubOptions)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "GitHub options were not provided to GitHub service").Err()
	}

	opts := req.GetGithubOptions()
	newRepo := &githubv3.Repository{
		Name:        strPtr(req.Name),
		Description: strPtr(req.Description),
		Visibility:  strPtr(strings.ToLower(opts.Parameters.Visibility.String())),
		AutoInit:    boolPtr(opts.AutoInit),
	}
	repo, _, err := s.rest.Repositories.Create(ctx, req.Owner, newRepo)
	if err != nil {
		return nil, err
	}

	resp := &sourcecontrolv1.CreateRepositoryResponse{
		Url: *repo.URL,
	}
	return resp, nil
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func (s *svc) CreatePullRequest(ctx context.Context, ref *RemoteRef, title, body string) (*PullRequestInfo, error) {
	req := &githubv3.NewPullRequest{
		Title:               strPtr(title),
		Head:                strPtr(ref.Ref),
		Base:                strPtr("master"),
		Body:                strPtr(body),
		MaintainerCanModify: boolPtr(true),
	}
	pr, _, err := s.rest.PullRequests.Create(ctx, ref.RepoOwner, ref.RepoName, req)
	if err != nil {
		return nil, err
	}
	return &PullRequestInfo{
		Number: pr.GetNumber(),
		// There are many possible URLs to return, but the HTML one is most human friendly
		HTMLURL: pr.GetHTMLURL(),
	}, nil
}

type CreateBranchRequest struct {
	// The base for the new branch.
	Ref *RemoteRef

	// The name of the new branch.
	BranchName string

	// Files and their content. Files will be clobbered with new content or created if they don't already exist.
	Files FileMap

	// The commit message for files added.
	CommitMessage string
}

// Creates a new branch with a commit containing files and pushes it to the remote.
func (s *svc) CreateBranch(ctx context.Context, req *CreateBranchRequest) error {
	cloneOpts := &git.CloneOptions{
		Depth:         1,
		URL:           fmt.Sprintf("https://github.com/%s/%s", req.Ref.RepoOwner, req.Ref.RepoName),
		ReferenceName: plumbing.NewBranchReferenceName(req.Ref.Ref),
		Auth:          s.rawAuth,
	}

	repo, err := git.CloneContext(ctx, memory.NewStorage(), memfs.New(), cloneOpts)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	checkoutOpts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(req.BranchName),
		Create: true,
	}
	if err := wt.Checkout(checkoutOpts); err != nil {
		return err
	}

	for filename, contents := range req.Files {
		fh, err := wt.Filesystem.Create(filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(fh, contents); err != nil {
			return err
		}
	}

	if err := wt.AddGlob("."); err != nil {
		return err
	}

	if _, err := wt.Commit(req.CommitMessage, &git.CommitOptions{}); err != nil {
		return err
	}

	pushOpts := &git.PushOptions{Auth: s.rawAuth}
	if err := repo.PushContext(ctx, pushOpts); err != nil {
		return err
	}

	return nil
}

func newService(config *githubv1.Config) Client {
	token := config.GetAccessToken()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	rest := githubv3.NewClient(httpClient)
	return &svc{
		graphQL: githubv4.NewClient(httpClient),
		rest: v3client{
			Repositories: rest.Repositories,
			PullRequests: rest.PullRequests,
			Issues:       rest.Issues,
		},
		rawAuth: &gittransport.BasicAuth{
			Username: "token",
			Password: token,
		},
	}
}

func (s *svc) GetFile(ctx context.Context, ref *RemoteRef, path string) (*File, error) {
	q := &getFileQuery{}
	params := map[string]interface{}{
		"owner":   githubv4.String(ref.RepoOwner),
		"name":    githubv4.String(ref.RepoName),
		"path":    githubv4.String(path),
		"ref":     githubv4.String(ref.Ref),
		"refPath": githubv4.String(fmt.Sprintf("%s:%s", ref.Ref, path)),
	}

	err := s.graphQL.Query(ctx, q, params)
	if err != nil {
		return nil, err
	}

	switch {
	case q.Repository.Ref.Commit.ID == nil:
		return nil, errors.New("ref not found")
	case q.Repository.Object.Blob.ID == nil:
		return nil, errors.New("object not found")
	case bool(q.Repository.Object.Blob.IsTruncated):
		return nil, errors.New("object was too large and was truncated by the API")
	case bool(q.Repository.Object.Blob.IsBinary):
		return nil, errors.New("object is a binary object and cannot be retrieved directly via the API")
	}

	f := &File{
		Path:     path,
		Contents: ioutil.NopCloser(strings.NewReader(string(q.Repository.Object.Blob.Text))),
		SHA:      string(q.Repository.Object.Blob.OID),
	}
	if len(q.Repository.Ref.Commit.History.Nodes) > 0 {
		f.LastModifiedTime = q.Repository.Ref.Commit.History.Nodes[0].CommittedDate.Time
		f.LastModifiedSHA = string(q.Repository.Ref.Commit.History.Nodes[0].OID)
	}

	return f, nil
}

func (s *svc) CompareCommits(ctx context.Context, ref *RemoteRef, compareSHA string) (*scgithubv1.CommitComparison, error) {
	comp, _, err := s.rest.Repositories.CompareCommits(ctx, ref.RepoOwner, ref.RepoName, compareSHA, ref.Ref)
	if err != nil {
		return nil, fmt.Errorf("Could not get compare status for %s and %s. %+v", ref.Ref, compareSHA, err)
	}

	status, ok := scgithubv1.CommitCompareStatus_value[strings.ToUpper(comp.GetStatus())]
	if !ok {
		return nil, fmt.Errorf("unknown status %s", comp.GetStatus())
	}

	return &scgithubv1.CommitComparison{
		Status: scgithubv1.CommitCompareStatus(status),
	}, nil
}

type Commit struct {
	Files []*githubv3.CommitFile
}

func (s *svc) GetCommit(ctx context.Context, ref *RemoteRef) (*Commit, error) {
	commit, _, err := s.rest.Repositories.GetCommit(ctx, ref.RepoOwner, ref.RepoName, ref.Ref)
	if err != nil {
		return nil, err
	}
	return &Commit{
		Files: commit.Files,
	}, nil
}
