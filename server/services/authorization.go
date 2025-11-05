package services

import "slices"

func reverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	reversedMap := make(map[V]K)
	for k, v := range m {
		reversedMap[v] = k
	}
	return reversedMap
}

type Role int

const (
	ROLE_NON_USER Role = iota
	ROLE_USER
	ROLE_ADMIN
)

var roleToStrings = map[Role]string{
	ROLE_NON_USER: "",
	ROLE_USER:     "user",
	ROLE_ADMIN:    "admin",
}

var stringToRoles = reverseMap(roleToStrings)

func (r Role) String() string {
	return roleToStrings[r]
}

func (r Role) Can(p Permission) bool {
	return slices.Contains(permissionsMap[r], p)
}

func (r Role) CanAny(perms ...Permission) bool {
	for _, perm := range perms {
		if r.Can(perm) {
			return true
		}
	}
	return false
}

func (r Role) CanAll(perms ...Permission) bool {
	for _, perm := range perms {
		if !r.Can(perm) {
			return false
		}
	}
	return true
}

func NewRole(roleString string) Role {
	role, ok := stringToRoles[roleString]
	if !ok {
		return ROLE_NON_USER
	}
	return role
}

type Permission int

const (
	PERMISSION_VIEW_WORD Permission = iota
	PERMISSION_ADD_WORD
	PERMISSION_DELETE_WORD
	PERMISSION_MODIFY_WORD
	PERMISSION_ADD_PROPOSAL
	PERMISSION_VIEW_ALL_PROPOSAL
	PERMISSION_VIEW_SELF_PROPOSAL
	PERMISSION_MODIFY_ALL_PROPOSAL
	PERMISSION_MODIFY_SELF_PROPOSAL
	PERMISSION_DELETE_ALL_PROPOSAL
	PERMISSION_DELETE_SELF_PROPOSAL
)

var permissionsMap = map[Role][]Permission{
	ROLE_ADMIN: {
		PERMISSION_VIEW_WORD,
		PERMISSION_ADD_WORD,
		PERMISSION_DELETE_WORD,
		PERMISSION_MODIFY_WORD,
		PERMISSION_ADD_PROPOSAL,
		PERMISSION_VIEW_ALL_PROPOSAL,
		PERMISSION_VIEW_SELF_PROPOSAL,
		PERMISSION_MODIFY_ALL_PROPOSAL,
		PERMISSION_MODIFY_SELF_PROPOSAL,
		PERMISSION_DELETE_ALL_PROPOSAL,
		PERMISSION_DELETE_SELF_PROPOSAL,
	},
	ROLE_USER: {
		PERMISSION_VIEW_WORD,
		PERMISSION_ADD_PROPOSAL,
		PERMISSION_VIEW_SELF_PROPOSAL,
		PERMISSION_MODIFY_SELF_PROPOSAL,
		PERMISSION_DELETE_SELF_PROPOSAL,
	},
	ROLE_NON_USER: {
		PERMISSION_VIEW_WORD,
	},
}
