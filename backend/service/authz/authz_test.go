package authz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	authzcfgv1 "github.com/lyft/clutch/backend/api/config/service/authz/v1"
)

func TestAssertPolicy(t *testing.T) {
	req := &authzv1.CheckRequest{
		Method:     "/clutch.v1.foo/Bar",
		ActionType: apiv1.ActionType_UPDATE,
		Resource:   "baz/qux",
	}

	tests := []struct {
		pol   *authzcfgv1.Policy
		match bool
	}{
		{
			pol:   &authzcfgv1.Policy{},
			match: true,
		},
		{
			pol: &authzcfgv1.Policy{
				ActionTypes: []apiv1.ActionType{apiv1.ActionType_READ, apiv1.ActionType_CREATE},
			},
			match: false,
		},
		{
			pol: &authzcfgv1.Policy{
				ActionTypes: []apiv1.ActionType{apiv1.ActionType_READ, apiv1.ActionType_UPDATE},
			},
			match: true,
		},
		{
			pol:   &authzcfgv1.Policy{Method: "/clutch.v1.foo/*"},
			match: true,
		},
		{
			pol:   &authzcfgv1.Policy{Method: "/*/Bar"},
			match: true,
		},
		{
			pol:   &authzcfgv1.Policy{Method: "/clutch.v1.foo/Bar", Resources: []string{"baz/*"}},
			match: true,
		},
		{
			pol:   &authzcfgv1.Policy{Method: "/clutch.v1.foo/Bar", Resources: []string{"baz/bang"}},
			match: false,
		},
		{
			pol:   &authzcfgv1.Policy{Method: "/clutch.v1.bar/Foo"},
			match: false,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			result := assertPolicy(tt.pol, req)
			assert.Equal(t, tt.match, result)
		})
	}
}

func newUserPrincipal(name string) *authzcfgv1.Principal {
	return &authzcfgv1.Principal{Type: &authzcfgv1.Principal_User{User: name}}
}

func newGroupPrincipal(name string) *authzcfgv1.Principal {
	return &authzcfgv1.Principal{Type: &authzcfgv1.Principal_Group{Group: name}}
}

func TestConfigToPrincipalRoleMap(t *testing.T) {
	toRoles := []string{"role-a", "role-b"}
	result := configToPrincipalRoleMap(&authzcfgv1.Config{
		RoleBindings: []*authzcfgv1.RoleBinding{
			{
				To: toRoles,
				Principals: []*authzcfgv1.Principal{
					newUserPrincipal("foo@example.com"),
					newUserPrincipal("bar@example.com"),
					newGroupPrincipal("group-z"),
				},
			},
		},
	})

	expected := principalToRoleMap{
		principalKey{principalType: user, name: "foo@example.com"}: toRoles,
		principalKey{principalType: user, name: "bar@example.com"}: toRoles,
		principalKey{principalType: group, name: "group-z"}:        toRoles,
	}

	assert.Equal(t, len(expected), len(result))
	for key, value := range result {
		roles, ok := expected[key]
		assert.True(t, ok)
		assert.ElementsMatch(t, roles, value)
	}
}

func TestConfigToRolePolicyMap(t *testing.T) {
	pol := []*authzcfgv1.Policy{{Method: "/foo/Bar"}}
	cfg := &authzcfgv1.Config{
		Roles: []*authzcfgv1.Role{
			{RoleName: "role-a", Policies: pol},
		},
	}
	result, err := configToRolePolicyMap(cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, roleToPolicyMap{"role-a": &authzcfgv1.Role{RoleName: "role-a", Policies: pol}}, result)
}
