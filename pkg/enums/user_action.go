package enums

type UserAction string

const (
	Register       UserAction = "register"
	Login                     = "login"
	Logout                    = "logout"
	VerifyEmail               = "verify_email"
	UpdateEmail               = "update_email"
	UpdatePassword            = "update_password"
)
