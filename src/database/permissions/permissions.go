package permissions

import (
	"wilin/src/database/roles"
)

type Permission string

const (
	VIEW_WORD   Permission = "view:word"
	ADD_WORD    Permission = "add:word"
	DELETE_WORD Permission = "delete:word"
	MODIFY_WORD Permission = "modify:word"
)

var adminPermissions = []Permission{
	VIEW_WORD,
	ADD_WORD,
	DELETE_WORD,
	MODIFY_WORD,
}

var userPermissions = []Permission{
	VIEW_WORD,
}

var nonUserPermissions = []Permission{
	VIEW_WORD,
}

func isPermissionInArray(permission Permission, permissionArray []Permission) bool {
	for _, permissionValue := range permissionArray {
		if permission == permissionValue {
			return true
		}
	}
	return false
}

func CanRolePermission(role roles.Role, permission Permission) bool {
	switch role {
	case roles.ADMIN:
		return isPermissionInArray(permission, adminPermissions)
	case roles.USER:
		return isPermissionInArray(permission, userPermissions)
	case roles.NON_USER:
		return isPermissionInArray(permission, nonUserPermissions)
	default:
		return false
	}
}
