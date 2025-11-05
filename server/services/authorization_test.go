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
	for _, testValue := range RoleStringValues {
		roleString := testValue.role.String()
		if strings.Compare(roleString, testValue.expected) != 0 {
			failTest(t, roleString, testValue.expected)
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
	for _, testValue := range RoleCanValues {
		output := testValue.role.Can(testValue.permission)
		if output != testValue.expected {
			failTest(t, output, testValue.expected)
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
