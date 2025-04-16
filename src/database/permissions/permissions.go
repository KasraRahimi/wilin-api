package permissions

import (
	"wilin/src/database/roles"
)

const (
	VIEW_WORD   = "view:word"
	ADD_WORD    = "add:word"
	DELETE_WORD = "delete:word"
	MODIFY_WORD = "modify:word"
)

var adminPermissions = []string{
	VIEW_WORD,
	ADD_WORD,
	DELETE_WORD,
	MODIFY_WORD,
}

var userPermissions = []string{
	VIEW_WORD,
}

var nonUserPermissions = []string{
	VIEW_WORD,
}

func isPermissionInArray(permission string, permissionArray []string) bool {
	for _, permissionValue := range permissionArray {
		if permission == permissionValue {
			return true
		}
	}
	return false
}

func CanRolePermission(role string, permission string) bool {
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
