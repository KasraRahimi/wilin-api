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
