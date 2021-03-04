package sourcecontrol

// <!-- START clutchdoc -->
// description: Creates repositories in a source control system.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	sourcecontrolconfigv1 "github.com/lyft/clutch/backend/api/config/module/sourcecontrol/v1"
	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/github"
)

const Name = "clutch.module.sourcecontrol"

func New(cfg *any.Any, _ *zap.Logger, _ tally.Scope) (module.Module, error) {
	config := &sourcecontrolconfigv1.Config{}
	if cfg != nil {
		if err := ptypes.UnmarshalAny(cfg, config); err != nil {
			return nil, err
		}
	}

	svc, ok := service.Registry["clutch.service.github"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := svc.(github.Client)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	return &mod{c, config}, nil
}

type mod struct {
	github github.Client
	config *sourcecontrolconfigv1.Config
}

func (m *mod) Register(r module.Registrar) error {
	sourcecontrolv1.RegisterSourceControlAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(sourcecontrolv1.RegisterSourceControlAPIHandler)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (m *mod) ManageableOrganizations(ctx context.Context) ([]*sourcecontrolv1.Entity, error) {
	organizations, err := m.github.ListOrganizations(ctx, github.CurrentUser)
	if err != nil {
		return nil, err
	}
	writableOrgs := []*sourcecontrolv1.Entity{}
	currentUser, err := m.github.GetUser(ctx, github.CurrentUser)
	if err != nil {
		return nil, err
	}
	if len(m.config.Owners) == 0 || contains(m.config.Owners, currentUser.GetLogin()) {
		writableOrgs = append(writableOrgs, &sourcecontrolv1.Entity{Name: currentUser.GetLogin(), PhotoUrl: currentUser.GetAvatarURL()})
	}
	for _, organization := range organizations {
		org, err := m.github.GetOrganization(ctx, organization.GetLogin())
		if err != nil {
			return nil, err
		}
		membership, err := m.github.GetOrgMembership(ctx, github.CurrentUser, org.GetLogin())
		if err != nil {
			return nil, err
		}
		if membership.GetRole() == "admin" || org.GetMembersCanCreatePrivateRepos() || org.GetMembersCanCreatePublicRepos() || org.GetMembersCanCreateRepos() {
			if len(m.config.Owners) == 0 || contains(m.config.Owners, org.GetLogin()) {
				writableOrgs = append(writableOrgs, &sourcecontrolv1.Entity{Name: org.GetLogin(), PhotoUrl: org.GetAvatarURL()})
			}
		}
	}

	sort.Slice(writableOrgs, func(i, j int) bool {
		return writableOrgs[i].Name < writableOrgs[j].Name
	})

	return writableOrgs, nil
}

func (m *mod) GetRepositoryOptions(ctx context.Context, req *sourcecontrolv1.GetRepositoryOptionsRequest) (*sourcecontrolv1.GetRepositoryOptionsResponse, error) {
	writableOrgs, err := m.ManageableOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	visibilityOptions := m.config.VisibilityOptions
	if visibilityOptions == nil {
		visibilityOptions = []sourcecontrolv1.Visibility{
			sourcecontrolv1.Visibility_PUBLIC,
			sourcecontrolv1.Visibility_PRIVATE,
		}
	}
	resp := &sourcecontrolv1.GetRepositoryOptionsResponse{
		AvailableOwners:   writableOrgs,
		VisibilityOptions: visibilityOptions,
	}
	return resp, nil
}

func (m *mod) CreateRepository(ctx context.Context, req *sourcecontrolv1.CreateRepositoryRequest) (*sourcecontrolv1.CreateRepositoryResponse, error) {
	desiredOwner := req.GetOwner()
	if len(m.config.Owners) != 0 && !contains(m.config.Owners, desiredOwner) {
		return nil, fmt.Errorf("cannot create repository under owner: %s", desiredOwner)
	}

	if len(m.config.VisibilityOptions) != 0 {
		opts := req.GetGithubOptions()
		desiredVisibility := opts.Parameters.Visibility

		visibilityOptions := make([]string, len(m.config.VisibilityOptions))
		for _, option := range m.config.VisibilityOptions {
			visibilityOptions = append(visibilityOptions, option.String())
		}
		if !contains(visibilityOptions, desiredVisibility.String()) {
			return nil, fmt.Errorf("cannot create repository with visibility: %s", desiredVisibility)
		}
	}
	return m.github.CreateRepository(ctx, req)
}
