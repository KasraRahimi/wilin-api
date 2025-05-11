package permissions

import (
	"wilin/src/database/roles"
)

type Permission string

const (
	VIEW_WORD          Permission = "view:word"
	ADD_WORD           Permission = "add:word"
	DELETE_WORD        Permission = "delete:word"
	MODIFY_WORD        Permission = "modify:word"
	ADD_PROPOSAL       Permission = "add:proposal"
	VIEW_ALL_PROPOSAL  Permission = "view:all:proposal"
	VIEW_SELF_PROPOSAL Permission = "view:self:proposal"
)

var permissionArray = map[roles.Role][]Permission{
	roles.ADMIN: {
		VIEW_WORD,
		ADD_WORD,
		DELETE_WORD,
		MODIFY_WORD,
		ADD_PROPOSAL,
		VIEW_ALL_PROPOSAL,
		VIEW_SELF_PROPOSAL,
	},
	roles.USER: {
		VIEW_WORD,
		ADD_PROPOSAL,
		VIEW_SELF_PROPOSAL,
	},
	roles.NON_USER: {
		VIEW_WORD,
	},
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
	return isPermissionInArray(permission, permissionArray[role])
}
