package services

type Role int

const (
	NonUser Role = iota
	User
	Admin
)

func (r Role) String() string {
	switch r {
	case NonUser:
		return ""
	case User:
		return "user"
	case Admin:
		return "admin"
	default:
		return ""
	}
}
