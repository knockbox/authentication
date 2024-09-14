package enums

var userRoleMap = map[UserRole]int{
	Banned:    0,
	Locked:    1,
	Pending:   2,
	User:      3,
	Moderator: 4,
	Admin:     5,
	Developer: 6,
}

type UserRole string

const (
	Banned    UserRole = "banned"
	Locked             = "locked"
	Pending            = "pending"
	User               = "user"
	Moderator          = "moderator"
	Admin              = "admin"
	Developer          = "developer"
)

func UserRoleFromString(role string) UserRole {
	switch role {
	case "banned":
		return Banned
	case "locked":
		return Locked
	case "pending":
		return Pending
	case "user":
		return User
	case "moderator":
		return Moderator
	case "admin":
		return Admin
	case "developer":
		return Developer
	default:
		return ""
	}
}

// IsForbidden determines if the role is forbidden
func (u UserRole) IsForbidden() bool {
	return u == Banned || u == Locked || u == Pending
}

// HasRequiredRole checks to see if the user has permission to access something
func (u UserRole) HasRequiredRole(required UserRole) bool {
	return userRoleMap[u] >= userRoleMap[required]
}

// IsDeveloperOrAdmin check to see if the user is a developer or admin
func (u UserRole) IsDeveloperOrAdmin() bool {
	return u == Developer || u == Admin
}
