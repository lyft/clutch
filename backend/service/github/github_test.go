package github

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	githubv3 "github.com/google/go-github/v32/github"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"

	githubv1 "github.com/lyft/clutch/backend/api/sourcecontrol/github/v1"
	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
)

const problem = "we've had a problem"

var (
	timestamp = time.Unix(1569010072, 0)
)

type getfileMock struct {
	v4client

	queryError        bool
	refID, objID      string
	truncated, binary bool
}

func (g *getfileMock) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	q, ok := query.(*getFileQuery)
	if !ok {
		panic("not a query")
	}

	if g.queryError {
		return errors.New(problem)
	}

	if g.refID != "" {
		q.Repository.Ref.Commit.ID = g.refID
		q.Repository.Ref.Commit.OID = githubv4.GitObjectID(g.refID)
	}
	if g.objID != "" {
		q.Repository.Object.Blob.OID = githubv4.GitObjectID(g.objID)
		q.Repository.Object.Blob.ID = g.objID
		q.Repository.Object.Blob.Text = "text"
	}

	q.Repository.Object.Blob.IsTruncated = githubv4.Boolean(g.truncated)
	q.Repository.Object.Blob.IsBinary = githubv4.Boolean(g.binary)

	q.Repository.Ref.Commit.History.Nodes = append(
		q.Repository.Ref.Commit.History.Nodes,
		struct {
			CommittedDate githubv4.DateTime
			OID           githubv4.GitObjectID
		}{githubv4.DateTime{Time: timestamp}, "otherSHA"},
	)
	return nil
}

var getFileTests = []struct {
	name    string
	v4      getfileMock
	errText string
}{
	{
		name:    "queryError",
		v4:      getfileMock{queryError: true},
		errText: problem,
	},
	{
		name:    "noRef",
		v4:      getfileMock{},
		errText: "ref not found",
	},
	{
		name:    "noObject",
		v4:      getfileMock{refID: "abcdef12345"},
		errText: "object not found",
	},
	{
		name:    "wasTruncated",
		v4:      getfileMock{refID: "abcdef12345", objID: "abcdef12345", truncated: true},
		errText: "truncated",
	},
	{
		name:    "wasBinary",
		v4:      getfileMock{refID: "abcdef12345", objID: "abcdef12345", binary: true},
		errText: "binary",
	},
	{
		name: "happyPath",
		v4:   getfileMock{refID: "abcdef12345", objID: "abcdef12345"},
	},
}

func TestGetFile(t *testing.T) {
	for _, tt := range getFileTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			s := &svc{graphQL: &tt.v4}
			f, err := s.GetFile(context.Background(),
				&RemoteRef{
					RepoOwner: "owner",
					RepoName:  "myRepo",
					Ref:       "master",
				},
				"data/foo",
			)

			if tt.errText != "" {
				a.Error(err)
				a.Contains(err.Error(), tt.errText)
				return
			}
			if err != nil {
				a.FailNow("unexpected error")
				return
			}

			contents, _ := ioutil.ReadAll(f.Contents)
			a.Equal("text", string(contents))
			a.Equal("data/foo", f.Path)
			a.Equal("abcdef12345", f.SHA)
			a.Equal("otherSHA", f.LastModifiedSHA)
			a.Equal(timestamp, f.LastModifiedTime)
		})
	}
}

type mockRepositories struct {
	actualOrg  string
	actualRepo *githubv3.Repository

	generalError           bool
	malformedEncodingError bool
}

func (m *mockRepositories) Create(ctx context.Context, org string, repo *githubv3.Repository) (*githubv3.Repository, *githubv3.Response, error) {
	m.actualRepo = repo
	m.actualOrg = org

	ret := &githubv3.Repository{
		URL: strPtr(fmt.Sprintf("https://example.com/%s/%s", org, *repo.Name)),
	}
	return ret, nil, nil
}

func (m *mockRepositories) GetContents(_ context.Context, _, _, _ string, _ *githubv3.RepositoryContentGetOptions) (*githubv3.RepositoryContent, []*githubv3.RepositoryContent, *githubv3.Response, error) {
	if m.generalError == true {
		return nil, nil, nil, errors.New(problem)
	}
	if m.malformedEncodingError {
		encoding := "unsupported"
		return &githubv3.RepositoryContent{Encoding: &encoding}, nil, nil, nil
	}
	return &githubv3.RepositoryContent{}, nil, nil, nil
}

func (m *mockRepositories) CompareCommits(ctx context.Context, owner, repo, base, head string) (*githubv3.CommitsComparison, *githubv3.Response, error) {
	if m.generalError {
		return nil, nil, errors.New(problem)
	}
	returnstr := "behind"
	return &githubv3.CommitsComparison{Status: &returnstr}, nil, nil
}

func (m *mockRepositories) GetCommit(ctx context.Context, owner, repo, sha string) (*githubv3.RepositoryCommit, *githubv3.Response, error) {
	file := "testfile.go"
	if m.generalError {
		return &githubv3.RepositoryCommit{}, &githubv3.Response{}, errors.New(problem)
	}
	return &githubv3.RepositoryCommit{
		Files: []*githubv3.CommitFile{
			{
				Filename: &file,
			},
		},
	}, &githubv3.Response{}, nil
}

var createRepoTests = []struct {
	req *sourcecontrolv1.CreateRepositoryRequest
}{
	{
		req: &sourcecontrolv1.CreateRepositoryRequest{
			Owner:       "foo",
			Name:        "bar",
			Description: "this is my repository",
			Options: &sourcecontrolv1.CreateRepositoryRequest_GithubOptions{GithubOptions: &githubv1.CreateRepositoryOptions{
				Parameters: &githubv1.RepositoryParameters{Visibility: githubv1.RepositoryParameters_PUBLIC},
			}},
		},
	},
	{
		req: &sourcecontrolv1.CreateRepositoryRequest{
			Owner:       "",
			Name:        "bar",
			Description: "this is also my repository",
			Options: &sourcecontrolv1.CreateRepositoryRequest_GithubOptions{GithubOptions: &githubv1.CreateRepositoryOptions{
				Parameters: &githubv1.RepositoryParameters{Visibility: githubv1.RepositoryParameters_PRIVATE},
			}},
		},
	},
}

func TestCreateRepository(t *testing.T) {
	for idx, tt := range createRepoTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			m := &mockRepositories{}
			s := &svc{rest: v3client{
				Repositories: m,
			}}

			resp, err := s.CreateRepository(context.Background(), tt.req)

			var expectedViz string
			switch tt.req.GetGithubOptions().Parameters.Visibility {
			case githubv1.RepositoryParameters_PUBLIC:
				expectedViz = "public"
			case githubv1.RepositoryParameters_PRIVATE:
				expectedViz = "private"
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.req.Owner, m.actualOrg)
			assert.Equal(t, tt.req.Name, *m.actualRepo.Name)
			assert.Equal(t, expectedViz, *m.actualRepo.Visibility)
			assert.Equal(t, tt.req.Description, *m.actualRepo.Description)
			assert.NotEmpty(t, resp.Url)
		})
	}
}

var compareCommitsTests = []struct {
	name         string
	errorText    string
	status       githubv1.CommitCompareStatus
	generalError bool
	mockRepo     *mockRepositories
}{
	{
		name:         "v3 error",
		generalError: true,
		errorText:    "Could not get compare status",
		mockRepo:     &mockRepositories{generalError: true},
	},
	{
		name:     "happy path",
		status:   githubv1.CommitCompareStatus_BEHIND,
		mockRepo: &mockRepositories{},
	},
}

func TestCompareCommits(t *testing.T) {
	t.Parallel()

	for _, tt := range compareCommitsTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)
			s := &svc{rest: v3client{
				Repositories: tt.mockRepo,
			}}

			comp, err := s.CompareCommits(
				context.Background(),
				&RemoteRef{
					RepoOwner: "owner",
					RepoName:  "myRepo",
					Ref:       "master",
				},
				"1234",
			)
			if tt.errorText != "" {
				a.Error(err)
				a.Contains(err.Error(), tt.errorText)
				return
			}
			if err != nil {
				a.FailNowf("unexpected error: %s", err.Error())
				return
			}
			a.Equal(comp.GetStatus(), tt.status)
			a.Nil(err)
		})
	}
}

var getCommitsTests = []struct {
	name      string
	errorText string
	mockRepo  *mockRepositories
	file      string
}{
	{
		name:      "v3 error",
		mockRepo:  &mockRepositories{generalError: true},
		errorText: "we've had a problem",
	},
	{
		name:     "happy path",
		mockRepo: &mockRepositories{},
		file:     "testfile.go",
	},
}

func TestGetCommit(t *testing.T) {
	t.Parallel()
	for _, tt := range getCommitsTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)
			s := &svc{rest: v3client{
				Repositories: tt.mockRepo,
			}}

			commit, err := s.GetCommit(context.Background(),
				&RemoteRef{
					RepoOwner: "owner",
					RepoName:  "myRepo",
					Ref:       "1234",
				},
			)

			if tt.errorText != "" {
				a.Error(err)
				a.Contains(err.Error(), tt.errorText)
				return
			}
			if err != nil {
				a.FailNowf("unexpected error: %s", err.Error())
				return
			}
			a.Equal(tt.file, *commit.Files[0].Filename)
			a.Nil(err)
		})
	}
}
