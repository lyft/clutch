package github

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	githubv3 "github.com/google/go-github/v35/github"
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
