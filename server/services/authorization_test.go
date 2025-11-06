package services_test

import (
	"strings"
	"testing"

	"wilin.info/api/server/services"
)

func failTest[T any](t *testing.T, got T, want T) {
	t.Errorf("got: %v, want: %v\n", got, want)
}

type RoleStringValue struct {
	role     services.Role
	expected string
}

var RoleStringValues = []RoleStringValue{
	{services.ROLE_NON_USER, ""},
	{services.ROLE_USER, "user"},
	{services.ROLE_ADMIN, "admin"},
}

func TestRoleString(t *testing.T) {
	for _, test := range RoleStringValues {
		roleString := test.role.String()
		if strings.Compare(roleString, test.expected) != 0 {
			failTest(t, roleString, test.expected)
		}
	}
}

type RoleCanValue struct {
	role       services.Role
	permission services.Permission
	expected   bool
}

var RoleCanValues = []RoleCanValue{
	{services.ROLE_NON_USER, services.PERMISSION_VIEW_WORD, true},
	{services.ROLE_NON_USER, services.PERMISSION_MODIFY_WORD, false},
	{services.ROLE_NON_USER, services.PERMISSION_DELETE_WORD, false},
	{services.ROLE_NON_USER, services.PERMISSION_VIEW_ALL_PROPOSAL, false},
	{services.ROLE_USER, services.PERMISSION_VIEW_SELF_PROPOSAL, true},
	{services.ROLE_USER, services.PERMISSION_VIEW_ALL_PROPOSAL, false},
	{services.ROLE_USER, services.PERMISSION_DELETE_WORD, false},
	{services.ROLE_USER, services.PERMISSION_MODIFY_SELF_PROPOSAL, true},
	{services.ROLE_ADMIN, services.PERMISSION_DELETE_ALL_PROPOSAL, true},
	{services.ROLE_ADMIN, services.PERMISSION_DELETE_WORD, true},
	{services.ROLE_ADMIN, services.PERMISSION_MODIFY_WORD, true},
	{services.ROLE_ADMIN, services.PERMISSION_VIEW_ALL_PROPOSAL, true},
}

func TestRoleCan(t *testing.T) {
	for _, test := range RoleCanValues {
		output := test.role.Can(test.permission)
		if output != test.expected {
			failTest(t, output, test.expected)
		}
	}
}

type RoleNewValue struct {
	input    string
	expected services.Role
}

var RoleNewValues = []RoleNewValue{
	{"admin", services.ROLE_ADMIN},
	{"Admin", services.ROLE_NON_USER},
	{"ADMIN", services.ROLE_NON_USER},
	{"aDmIn", services.ROLE_NON_USER},
	{"user", services.ROLE_USER},
	{"User", services.ROLE_NON_USER},
	{"USER", services.ROLE_NON_USER},
	{"nada", services.ROLE_NON_USER},
	{"non user", services.ROLE_NON_USER},
	{"guest", services.ROLE_NON_USER},
}

func TestRoleNew(t *testing.T) {
	for _, test := range RoleNewValues {
		role := services.NewRole(test.input)
		if role != test.expected {
			failTest(t, role, test.expected)
		}
	}
}

type RoleCanArrValue struct {
	role        services.Role
	permissions []services.Permission
	expected    bool
}

var roleCanAnyValues = []RoleCanArrValue{
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		true,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		false,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		false,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_DELETE_WORD,
		},
		true,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_DELETE_ALL_PROPOSAL,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		false,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		true,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_ADD_WORD,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_VIEW_ALL_PROPOSAL,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
		},
		true,
	},
}

func TestRoleCanAny(t *testing.T) {
	for _, test := range roleCanAnyValues {
		result := test.role.CanAny(test.permissions...)
		if result != test.expected {
			failTest(t, result, test.expected)
		}
	}
}

var roleCanAllValues = []RoleCanArrValue{
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
		},
		true,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		false,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		false,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		false,
	},
	{
		services.ROLE_NON_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		false,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_DELETE_WORD,
		},
		false,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
		},
		false,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_DELETE_ALL_PROPOSAL,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		false,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_USER,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
		},
		false,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_MODIFY_WORD,
		},
		true,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
		},
		true,
	},
	{
		services.ROLE_ADMIN,
		[]services.Permission{
			services.PERMISSION_VIEW_WORD,
			services.PERMISSION_ADD_WORD,
			services.PERMISSION_DELETE_WORD,
			services.PERMISSION_MODIFY_WORD,
			services.PERMISSION_ADD_PROPOSAL,
			services.PERMISSION_VIEW_ALL_PROPOSAL,
			services.PERMISSION_VIEW_SELF_PROPOSAL,
			services.PERMISSION_MODIFY_ALL_PROPOSAL,
			services.PERMISSION_MODIFY_SELF_PROPOSAL,
			services.PERMISSION_DELETE_ALL_PROPOSAL,
			services.PERMISSION_DELETE_SELF_PROPOSAL,
		},
		true,
	},
}

func TestRoleCanAll(t *testing.T) {
	for _, test := range roleCanAllValues {
		result := test.role.CanAll(test.permissions...)
		if result != test.expected {
			failTest(t, result, test.expected)
		}
	}
}
