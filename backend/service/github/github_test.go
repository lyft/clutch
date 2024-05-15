package github

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	githubv3 "github.com/google/go-github/v54/github"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	githubconfigv1 "github.com/lyft/clutch/backend/api/config/service/github/v1"
	githubv1 "github.com/lyft/clutch/backend/api/sourcecontrol/github/v1"
	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/service/authn"
)

const problem = "we've had a problem"

var timestamp = time.Unix(1569010072, 0)

func intPtr(i int) *int {
	return &i
}

type getfileMock struct {
	v4client

	queryError        bool
	refID, objID      string
	truncated, binary bool
}

type getfilepathMock struct {
	v4client

	queryError   bool
	refID, objID string
	entries      []*FileEntry
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

func (g *getfilepathMock) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	q, ok := query.(*getFilePathQuery)
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
		q.Repository.Object.Tree.Entries = append(
			q.Repository.Object.Tree.Entries,
			struct {
				Name githubv4.String
				Type githubv4.String
			}{Name: githubv4.String(g.entries[0].Name), Type: githubv4.String(g.entries[0].Type)})
	}

	q.Repository.Ref.Commit.History.Nodes = append(
		q.Repository.Ref.Commit.History.Nodes,
		struct {
			CommittedDate githubv4.DateTime
			OID           githubv4.GitObjectID
		}{githubv4.DateTime{Time: timestamp}, "otherSHA"},
	)

	return nil
}

type mockRoundTripper struct {
	http.RoundTripper
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if m.err != nil {
		return m.resp, m.err
	}
	return m.resp, nil
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name          string
		headerName    string
		headerValue   string
		nilResp       bool
		respErr       error
		expectedGauge float64
	}{
		{
			name:          "gauges rate limit remaining",
			headerName:    "X-RateLimit-Remaining",
			headerValue:   "10",
			expectedGauge: 10,
		},
		{
			name:        "only updates gauge when header present",
			headerName:  "x-unrelated-header",
			headerValue: "foobar",
		},
		{
			name:        "only updates gauge when header value valid type",
			headerName:  "X-RateLimit-Remaining",
			headerValue: "non-int",
		},
		{
			name:    "round trip missing response update gauge",
			respErr: errors.New("upstream error"),
		},
		{
			name:          "round trip error with response updates gauge",
			respErr:       errors.New("upstream error"),
			headerName:    "X-RateLimit-Remaining",
			headerValue:   "10",
			expectedGauge: 10,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scope := tally.NewTestScope("", nil)
			h := http.Header{}
			h.Set(tt.headerName, tt.headerValue)
			r := &http.Response{Header: h}
			if tt.nilResp {
				r = nil
			}
			mockRT := &mockRoundTripper{resp: r, err: tt.respErr}
			st := &StatsRoundTripper{
				Wrapped: mockRT,
				scope:   scope,
			}

			resp, err := st.RoundTrip(&http.Request{})

			assert.Equal(t, r, resp)
			if tt.respErr == nil {
				assert.Nil(t, err)
			}

			g := scope.Snapshot().Gauges()["rate_limit_remaining+"]
			if tt.expectedGauge != 0 {
				assert.NotNil(t, g)
				assert.Equal(t, tt.expectedGauge, g.Value())
			} else {
				assert.Nil(t, g)
			}
		})
	}
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

			contents, _ := io.ReadAll(f.Contents)
			a.Equal("text", string(contents))
			a.Equal("data/foo", f.Path)
			a.Equal("abcdef12345", f.SHA)
			a.Equal("otherSHA", f.LastModifiedSHA)
			a.Equal(timestamp, f.LastModifiedTime)
		})
	}
}

var filePathEntries = []*FileEntry{{Name: "foo", Type: "blob"}}

var getFilePathTests = []struct {
	name    string
	v4      getfilepathMock
	errText string
}{
	{
		name:    "queryError",
		v4:      getfilepathMock{queryError: true},
		errText: problem,
	},
	{
		name:    "noRef",
		v4:      getfilepathMock{},
		errText: "ref not found",
	},
	{
		name:    "noObject",
		v4:      getfilepathMock{refID: "abcdef12345"},
		errText: "path not found",
	},
	{
		name: "happyPath",
		v4:   getfilepathMock{refID: "abcdef12345", objID: "abcdef12345", entries: filePathEntries},
	},
}

func TestGetFilePath(t *testing.T) {
	for _, tt := range getFilePathTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			s := &svc{graphQL: &tt.v4}
			f, err := s.GetFilePath(context.Background(),
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

			a.Equal("data/foo", f.Path)
			a.Equal(filePathEntries, filePathEntries)
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
		HTMLURL: strPtr(fmt.Sprintf("https://example.com/%s/%s", org, *repo.Name)),
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

func (m *mockRepositories) CompareCommits(ctx context.Context, owner, repo, base, head string, opts *githubv3.ListOptions) (*githubv3.CommitsComparison, *githubv3.Response, error) {
	if m.generalError {
		return nil, nil, errors.New(problem)
	}
	returnstr := "behind"
	shaStr := "astdfsaohecra"
	return &githubv3.CommitsComparison{
		Status: &returnstr,
		Commits: []*githubv3.RepositoryCommit{
			{SHA: &shaStr},
		},
	}, nil, nil
}

func (m *mockRepositories) GetCommit(ctx context.Context, owner, repo, sha string, opts *githubv3.ListOptions) (*githubv3.RepositoryCommit, *githubv3.Response, error) {
	file := "testfile.go"
	message := "committing some changes (#1)"
	authorLogin := "foobar"
	if m.generalError {
		return &githubv3.RepositoryCommit{}, &githubv3.Response{}, errors.New(problem)
	}
	return &githubv3.RepositoryCommit{
		Files: []*githubv3.CommitFile{
			{
				Filename: &file,
			},
		},
		Commit: &githubv3.Commit{
			Message: &message,
			Author: &githubv3.CommitAuthor{
				Login: &authorLogin,
			},
		},
	}, nil, nil
}

type mockUsers struct {
	user        githubv3.User
	defaultUser string
}

func (m *mockUsers) Get(ctx context.Context, user string) (*githubv3.User, *githubv3.Response, error) {
	var login string
	if login = user; user == "" {
		login = m.defaultUser
	}
	ret := &githubv3.User{
		Login: &login,
	}
	m.user = *ret
	return ret, nil, nil
}

var createRepoTests = []struct {
	req   *sourcecontrolv1.CreateRepositoryRequest
	users *mockUsers
}{
	{
		req: &sourcecontrolv1.CreateRepositoryRequest{
			Owner:       "organization",
			Name:        "bar",
			Description: "this is an org repository",
			Options: &sourcecontrolv1.CreateRepositoryRequest_GithubOptions{GithubOptions: &githubv1.CreateRepositoryOptions{
				Parameters: &githubv1.RepositoryParameters{Visibility: githubv1.RepositoryParameters_PUBLIC},
			}},
		},
		users: &mockUsers{},
	},
	{
		req: &sourcecontrolv1.CreateRepositoryRequest{
			Owner:       "user",
			Name:        "bar",
			Description: "this is my repository",
			Options: &sourcecontrolv1.CreateRepositoryRequest_GithubOptions{GithubOptions: &githubv1.CreateRepositoryOptions{
				Parameters: &githubv1.RepositoryParameters{Visibility: githubv1.RepositoryParameters_PRIVATE},
			}},
		},
		users: &mockUsers{
			defaultUser: "user",
		},
	},
}

func TestCreateRepository(t *testing.T) {
	for idx, tt := range createRepoTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			r := &mockRepositories{}
			s := &svc{rest: v3client{
				Repositories: r,
				Users:        tt.users,
			}}

			resp, err := s.CreateRepository(context.Background(), tt.req)

			var expectedOwner string
			if expectedOwner = tt.req.Owner; tt.req.Owner == "user" {
				expectedOwner = ""
			}

			var expectedPrivate bool
			switch tt.req.GetGithubOptions().Parameters.Visibility {
			case githubv1.RepositoryParameters_PUBLIC:
				expectedPrivate = false
			case githubv1.RepositoryParameters_PRIVATE:
				expectedPrivate = true
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedOwner, r.actualOrg)
			assert.Equal(t, tt.req.Name, *r.actualRepo.Name)
			assert.Equal(t, expectedPrivate, *r.actualRepo.Private)
			assert.Equal(t, tt.req.Description, *r.actualRepo.Description)
			assert.NotEmpty(t, resp.Url)
		})
	}
}

var getUserTests = []struct {
	username string
}{
	{
		username: "foobar",
	},
}

func TestGetUser(t *testing.T) {
	for idx, tt := range getUserTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			u := &mockUsers{}
			s := &svc{rest: v3client{
				Users: u,
			}}

			resp, err := s.GetUser(context.Background(), tt.username)

			assert.NoError(t, err)
			assert.Equal(t, u.user.GetLogin(), resp.GetLogin())
		})
	}
}

var compareCommitsTests = []struct {
	name         string
	errorText    string
	status       string
	generalError bool
	mockRepo     *mockRepositories
}{
	{
		name:         "v3 error",
		generalError: true,
		errorText:    "could not get comparison",
		mockRepo:     &mockRepositories{generalError: true},
	},
	{
		name:     "happy path",
		status:   "behind",
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
			a.NotNil(comp.Commits)
			a.Nil(err)
		})
	}
}

var getCommitsTests = []struct {
	name            string
	errorText       string
	mockRepo        *mockRepositories
	file            string
	message         string
	authorLogin     string
	authorAvatarURL string
	authorID        int64
	parentRef       string
}{
	{
		name:      "v3 error",
		mockRepo:  &mockRepositories{generalError: true},
		errorText: "we've had a problem",
	},
	{
		name:            "happy path",
		mockRepo:        &mockRepositories{},
		file:            "testfile.go",
		message:         "committing some changes (#1)",
		authorAvatarURL: "https://foo.bar/baz.png",
		authorID:        1234,
		parentRef:       "test",
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
			a.Equal(tt.message, commit.Message)
			if commit.Author != nil {
				a.Equal(tt.authorAvatarURL, *commit.Author.AvatarURL)
				a.Equal(tt.authorID, *commit.Author.ID)
			}
			if commit.ParentRef != "" {
				a.Equal(tt.parentRef, commit.ParentRef)
			}
			a.Nil(err)
		})
	}
}

type mockOrganizations struct {
	actualOrg  string
	actualUser string

	generalError bool
	authError    bool
}

func (m *mockOrganizations) Get(ctx context.Context, org string) (*githubv3.Organization, *githubv3.Response, error) {
	m.actualOrg = org
	if m.generalError == true {
		return nil, nil, errors.New(problem)
	}
	return &githubv3.Organization{
		Name: &org,
	}, nil, nil
}

func (m *mockOrganizations) List(ctx context.Context, user string, opts *githubv3.ListOptions) ([]*githubv3.Organization, *githubv3.Response, error) {
	m.actualUser = user
	if m.generalError == true {
		return nil, nil, errors.New(problem)
	}
	return []*githubv3.Organization{}, nil, nil
}

func (m *mockOrganizations) GetOrgMembership(ctx context.Context, user, org string) (*githubv3.Membership, *githubv3.Response, error) {
	m.actualOrg = org
	m.actualUser = user
	if m.generalError {
		return nil, &githubv3.Response{Response: &http.Response{StatusCode: 500}}, errors.New(problem)
	}
	if m.authError {
		return nil, &githubv3.Response{Response: &http.Response{StatusCode: 403}}, nil
	}
	return &githubv3.Membership{}, nil, nil
}

var getOrganizationTests = []struct {
	name      string
	errorText string
	mockOrgs  *mockOrganizations
	org       string
}{
	{
		name:      "v3 error",
		mockOrgs:  &mockOrganizations{generalError: true},
		errorText: "we've had a problem",
		org:       "testing",
	},
	{
		name:     "v3 error",
		mockOrgs: &mockOrganizations{authError: true},
		org:      "testing",
	},
	{
		name:     "happy path",
		mockOrgs: &mockOrganizations{},
		org:      "testing",
	},
}

func TestGetOrganization(t *testing.T) {
	for idx, tt := range getOrganizationTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			s := &svc{rest: v3client{
				Organizations: tt.mockOrgs,
			}}

			resp, err := s.GetOrganization(context.Background(), tt.org)

			if tt.errorText != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, resp.GetName(), tt.org)
				assert.Equal(t, tt.mockOrgs.actualOrg, tt.org)
			}
		})
	}
}

var listOrganizationsTests = []struct {
	name      string
	errorText string
	mockOrgs  *mockOrganizations
	username  string
}{
	{
		name:      "v3 error",
		mockOrgs:  &mockOrganizations{generalError: true},
		errorText: "we've had a problem",
		username:  "foobar",
	},
	{
		name:     "happy path",
		mockOrgs: &mockOrganizations{},
		username: "foobar",
	},
}

func TestListOrganizations(t *testing.T) {
	for idx, tt := range listOrganizationsTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			s := &svc{rest: v3client{
				Organizations: tt.mockOrgs,
			}}

			resp, err := s.ListOrganizations(context.Background(), tt.username)

			if tt.errorText != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(resp), 0)
				assert.Equal(t, tt.mockOrgs.actualUser, tt.username)
			}
		})
	}
}

var getOrgMembershipTests = []struct {
	name      string
	errorText string
	mockOrgs  *mockOrganizations
	username  string
	org       string
}{
	{
		name:      "v3 error",
		mockOrgs:  &mockOrganizations{generalError: true},
		errorText: "we've had a problem",
		username:  "foobar",
		org:       "testing",
	},
	{
		name:     "happy path",
		mockOrgs: &mockOrganizations{},
		username: "foobar",
		org:      "testing",
	},
}

func TestGetOrgMembership(t *testing.T) {
	for idx, tt := range getOrgMembershipTests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			s := &svc{rest: v3client{
				Organizations: tt.mockOrgs,
			}}

			_, err := s.GetOrgMembership(context.Background(), tt.username, tt.org)

			if tt.errorText != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockOrgs.actualOrg, tt.org)
				assert.Equal(t, tt.mockOrgs.actualUser, tt.username)
			}
		})
	}
}

type getRepositoryMock struct {
	v4client
	branchName string
}

func (g *getRepositoryMock) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	q, ok := query.(*getRepositoryQuery)
	if !ok {
		panic("not a query")
	}
	q.Repository.DefaultBranchRef.Name = g.branchName
	return nil
}

var getDefaultBranchTests = []struct {
	name string
	v4   getRepositoryMock

	wantDefaultBranch string
}{
	{
		name:              "1. default repo with main branch",
		v4:                getRepositoryMock{branchName: "main"},
		wantDefaultBranch: "main",
	},
	{
		name:              "2. default repo with master branch",
		v4:                getRepositoryMock{branchName: "master"},
		wantDefaultBranch: "master",
	},
}

func TestGetRepository(t *testing.T) {
	t.Parallel()
	for _, tt := range getDefaultBranchTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)
			s := &svc{graphQL: &tt.v4}

			repo, err := s.GetRepository(context.Background(),
				&RemoteRef{
					RepoOwner: "owner",
					RepoName:  "myRepo",
				},
			)
			if err != nil {
				a.FailNowf("unexpected error: %s", err.Error())
				return
			}

			gotDefaultBranch := repo.DefaultBranch

			a.Equal(gotDefaultBranch, tt.wantDefaultBranch)
			a.Nil(err)
		})
	}
}

type mockPullRequests struct {
	generalError bool

	actualNumber     int
	actualHTMLURL    string
	actualBranchName string
}

// Dummy mock of Create API so mockPullRequests implements v3pullrequests
func (m *mockPullRequests) Create(ctx context.Context, owner string, repo string, pull *githubv3.NewPullRequest) (*githubv3.PullRequest, *githubv3.Response, error) {
	return &githubv3.PullRequest{}, &githubv3.Response{}, nil
}

// Mock of ListPullRequestsWithCommit API
func (m *mockPullRequests) ListPullRequestsWithCommit(ctx context.Context, owner, repo, sha string, opts *githubv3.ListOptions) ([]*githubv3.PullRequest, *githubv3.Response, error) {
	if m.generalError {
		return nil, nil, errors.New(problem)
	}

	m.actualNumber = 1347
	m.actualHTMLURL = fmt.Sprintf("https://github.com/%s/%s/pull/%s", owner, repo, strconv.Itoa(m.actualNumber))
	m.actualBranchName = "my-branch"

	return []*githubv3.PullRequest{
		{
			Number:  intPtr(m.actualNumber),
			HTMLURL: strPtr(m.actualHTMLURL),
			Head: &githubv3.PullRequestBranch{
				Ref: strPtr(m.actualBranchName),
				SHA: strPtr(sha),
				Repo: &githubv3.Repository{
					Name: strPtr(repo),
				},
				User: &githubv3.User{
					Login: strPtr("octocat"),
				},
			},
		},
	}, nil, nil
}

var listPullRequestsWithCommitTests = []struct {
	name        string
	errorText   string
	mockPullReq *mockPullRequests
	repoOwner   string
	repoName    string
	ref         string
	sha         string
}{
	{
		name:        "happy path",
		mockPullReq: &mockPullRequests{},
		repoOwner:   "my-org",
		repoName:    "my-repo",
		ref:         "my-branch",
		sha:         "asdf12345",
	},
	{
		name:        "v3 client error",
		mockPullReq: &mockPullRequests{generalError: true},
		errorText:   "we've had a problem",
		repoOwner:   "my-org",
		repoName:    "my-repo",
		ref:         "my-branch",
		sha:         "asdf12345",
	},
}

func TestListPullRequestsWithCommit(t *testing.T) {
	for _, tt := range listPullRequestsWithCommitTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &svc{rest: v3client{
				PullRequests: tt.mockPullReq,
			}}

			resp, err := s.ListPullRequestsWithCommit(
				context.Background(),
				&RemoteRef{
					RepoOwner: tt.repoOwner,
					RepoName:  tt.repoName,
					Ref:       tt.ref,
				},
				tt.sha, nil)

			if tt.errorText != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorText)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(resp))
				assert.Equal(t, tt.mockPullReq.actualNumber, resp[0].Number)
				assert.Equal(t, tt.mockPullReq.actualHTMLURL, resp[0].HTMLURL)
				assert.Equal(t, tt.mockPullReq.actualBranchName, resp[0].BranchName)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	cfg := &githubconfigv1.Config{}
	_, err := newService(cfg, tally.NoopScope, zap.NewNop())
	assert.Error(t, err)

	cfg.Auth = &githubconfigv1.Config_AccessToken{AccessToken: "aaa"}
	s, err := newService(cfg, tally.NoopScope, zap.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.Nil(t, s.(*svc).appTransport)
	assert.Equal(t, s.(*svc).personalAccessToken, "aaa")

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	kk, _ := x509.MarshalPKCS8PrivateKey(pk)

	cfg.Auth = &githubconfigv1.Config_AppConfig{
		AppConfig: &githubconfigv1.AppConfig{
			AppId:          123456,
			InstallationId: 789,
			Pem: &githubconfigv1.AppConfig_KeyPem{
				KeyPem: string(pem.EncodeToMemory(&pem.Block{Type: "RSA_PRIVATE_KEY", Bytes: kk})),
			},
		},
	}
	s, err = newService(cfg, tally.NoopScope, zap.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.Empty(t, s.(*svc).personalAccessToken)
	assert.NotNil(t, s.(*svc).appTransport)
}

func TestCommitAuthorFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = authn.ContextWithAnonymousClaims(ctx)
	commitTime := time.Now()

	result := commitOptionsFromClaims(ctx, commitTime)
	assert.Equal(t, "Anonymous User via Clutch", result.Author.Name)
	assert.Equal(t, "", result.Author.Email)
	assert.Equal(t, commitTime, result.Author.When)
	assert.Equal(t, true, result.All)

	ctx = authn.ContextWithClaims(ctx, &authn.Claims{
		StandardClaims: &jwt.StandardClaims{
			Subject: "daniel@example.com",
		},
	})

	commitTime = time.Now()
	result = commitOptionsFromClaims(ctx, commitTime)
	assert.Equal(t, "daniel@example.com via Clutch", result.Author.Name)
	assert.Equal(t, "daniel@example.com", result.Author.Email)
	assert.Equal(t, commitTime, result.Author.When)
	assert.Equal(t, true, result.All)

	ctx = authn.ContextWithClaims(ctx, &authn.Claims{
		StandardClaims: &jwt.StandardClaims{
			Subject: "daniel123",
		},
	})

	commitTime = time.Now()
	result = commitOptionsFromClaims(ctx, commitTime)
	assert.Equal(t, "daniel123 via Clutch", result.Author.Name)
	assert.Equal(t, "", result.Author.Email)
	assert.Equal(t, commitTime, result.Author.When)
	assert.Equal(t, true, result.All)
}
