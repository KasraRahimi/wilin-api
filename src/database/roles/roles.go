package roles

type Role string

const (
	NON_USER Role = "nonUser"
	USER     Role = "user"
	ADMIN    Role = "admin"
)
