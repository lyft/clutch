package authz

// <!-- START clutchdoc -->
// description: Evaluates the configured RBAC policies against user and resource pairs.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	authzcfgv1 "github.com/lyft/clutch/backend/api/config/service/authz/v1"
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.authz"

type Client interface {
	Check(ctx context.Context, request *authzv1.CheckRequest) (*authzv1.CheckResponse, error)
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &authzcfgv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}
	return newStaticImpl(logger, config)
}

// Private unexported struct that is used to index a map of principals to their roles.
// The proto struct is not used directly because it has internal fields that could interfere with
// hashing equally.
type principalKey struct {
	name          string
	principalType principalKeyType
}

type principalKeyType int

const (
	user principalKeyType = iota
	group
)

type principalToRoleMap map[principalKey][]string
type roleToPolicyMap map[string]*authzcfgv1.Role

type staticImpl struct {
	// Map of policy role names to the policy object.
	roleToPolicy roleToPolicyMap

	// Map of principals (i.e users or groups) to a list of role names that the principal is assigned.
	principalToRole principalToRoleMap
}

func configToRolePolicyMap(config *authzcfgv1.Config) (roleToPolicyMap, error) {
	// Pre-compute a role to policy map for faster lookup.
	roleToPolicy := make(roleToPolicyMap, len(config.Roles))
	for _, role := range config.Roles {
		if _, ok := roleToPolicy[role.RoleName]; ok {
			return nil, fmt.Errorf("duplicate role '%s'", role.RoleName)
		}
		roleToPolicy[role.RoleName] = role
	}

	// Verify all bound roles exist.
	for _, rb := range config.RoleBindings {
		for _, roleName := range rb.To {
			if _, ok := roleToPolicy[roleName]; !ok {
				return nil, fmt.Errorf("attempted to bind to non-existent role '%s'", roleName)
			}
		}
	}
	return roleToPolicy, nil
}

func configToPrincipalRoleMap(config *authzcfgv1.Config) principalToRoleMap {
	// Pre-compute a principal to role map.
	principalToRole := make(principalToRoleMap)
	for _, rb := range config.RoleBindings {
		for _, p := range rb.Principals {
			var key principalKey
			switch p.Type.(type) {
			case *authzcfgv1.Principal_User:
				key = principalKey{
					name:          p.GetUser(),
					principalType: user,
				}
			case *authzcfgv1.Principal_Group:
				key = principalKey{
					name:          p.GetGroup(),
					principalType: group,
				}
			}
			principalToRole[key] = append(principalToRole[key], rb.To...)
		}
	}

	// Unique the role lists for each principal.
	for k, v := range principalToRole {
		vset := make(map[string]struct{})
		for _, s := range v {
			vset[s] = struct{}{}
		}
		vv := make([]string, 0, len(vset))
		for s := range vset {
			vv = append(vv, s)
		}
		principalToRole[k] = vv
	}

	return principalToRole
}

func newStaticImpl(logger *zap.Logger, config *authzcfgv1.Config) (Client, error) {
	// Compute map of role to policy.
	roleToPolicy, err := configToRolePolicyMap(config)
	if err != nil {
		return nil, err
	}

	// Compute map of principal (user or group) to list of roles.
	principalToRole := configToPrincipalRoleMap(config)

	// Save on the struct for lookup at runtime.
	return &staticImpl{
		principalToRole: principalToRole,
		roleToPolicy:    roleToPolicy,
	}, nil
}

func assertPolicy(pol *authzcfgv1.Policy, req *authzv1.CheckRequest) bool {
	// ActionTypes: if none specified or request action is in the policy's list, OK.
	if len(pol.ActionTypes) != 0 {
		actionMatch := false
		for _, at := range pol.ActionTypes {
			if req.ActionType == at {
				actionMatch = true
				break
			}
		}
		if !actionMatch {
			return false
		}
	}

	// Method: if none specified or match, OK.
	if pol.Method != "" && !middleware.MatchMethodOrResource(pol.Method, req.Method) {
		return false
	}

	// Resources: if none specified or match, OK.
	if len(pol.Resources) != 0 {
		resMatch := false
		for _, r := range pol.Resources {
			if middleware.MatchMethodOrResource(r, req.Resource) {
				resMatch = true
				break
			}
		}
		if !resMatch {
			return false
		}
	}

	// No problem, OK.
	return true
}

func (s *staticImpl) evaluate(roles []string, req *authzv1.CheckRequest) *authzv1.CheckResponse {
	if len(roles) == 0 {
		return &authzv1.CheckResponse{
			Decision: authzv1.Decision_DENY,
		}
	}

	for _, roleName := range roles {
		role := s.roleToPolicy[roleName]
		for _, policy := range role.Policies {
			if assertPolicy(policy, req) {
				return &authzv1.CheckResponse{
					Decision: authzv1.Decision_ALLOW,
				}
			}
		}
	}

	return &authzv1.CheckResponse{
		Decision: authzv1.Decision_DENY,
	}
}

func (s *staticImpl) Check(ctx context.Context, req *authzv1.CheckRequest) (*authzv1.CheckResponse, error) {
	// Gather all of the roles for the user and/or groups.
	var roles []string
	if req.Subject.User != "" {
		roles = s.principalToRole[principalKey{
			name:          req.Subject.User,
			principalType: user,
		}]
	}

	for _, g := range req.Subject.Groups {
		if g == "" {
			continue
		}
		k := principalKey{name: g, principalType: group}
		roles = append(roles, s.principalToRole[k]...)
	}

	// Evaluate!
	resp := s.evaluate(roles, req)
	return resp, nil
}
