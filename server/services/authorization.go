package services

type Role int

const (
	RoleNonUser Role = iota
	RoleUser
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleNonUser:
		return ""
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return ""
	}
}
