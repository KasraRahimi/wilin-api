package services

type Role int

const (
	ROLE_NON_USER Role = iota
	ROLE_USER
	ROLE_ADMIN
)

func (r Role) String() string {
	switch r {
	case ROLE_NON_USER:
		return ""
	case ROLE_USER:
		return "user"
	case ROLE_ADMIN:
		return "admin"
	default:
		return ""
	}
}

func (r Role) Can(p Permission) bool {
	return p.IsPermissionInArray(permissionsMap[r])
}

func NewRole(role string) Role {
	switch role {
	case ROLE_NON_USER.String():
		return ROLE_NON_USER
	case ROLE_USER.String():
		return ROLE_USER
	case ROLE_ADMIN.String():
		return ROLE_ADMIN
	default:
		return ROLE_NON_USER
	}
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

func (p Permission) IsPermissionInArray(permissionArray []Permission) bool {
	for _, permissionValue := range permissionArray {
		if p == permissionValue {
			return true
		}
	}
	return false
}
