package models

import (
	"github.com/google/uuid"
	"github.com/knockbox/authentication/pkg/enums"
	"github.com/knockbox/authentication/pkg/payloads"
	"github.com/knockbox/authentication/pkg/utils"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

// User define a User in our database.
type User struct {
	Id        uint           `db:"id"`
	AccountId uuid.UUID      `db:"account_id"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	Email     string         `db:"email"`
	Role      enums.UserRole `db:"role"`
}

// NewUser creates a new User with an auto-generated uuid.UUID and role set to enums.User.
func NewUser() *User {
	return &User{
		Id:        0,
		AccountId: uuid.New(),
		Username:  "",
		Password:  "",
		Email:     "",
		Role:      enums.User,
	}
}

// ApplyRegister updates the user with the values from the payloads.UserRegister.
func (u *User) ApplyRegister(payload *payloads.UserRegister) error {
	u.Username = payload.Username
	u.Email = payload.Email

	pwd, err := utils.GeneratePassword(payload.Password)
	if err != nil {
		return err
	}
	u.Password = pwd

	return nil
}

// ApplyUpdate updates the user with the values from the payloads.UserUpdate.
func (u *User) ApplyUpdate(payload *payloads.UserUpdate) error {
	if payload.Email != nil {
		u.Email = *payload.Email
	}

	if payload.Password != nil {
		pwd, err := utils.GeneratePassword(*payload.Password)
		if err != nil {
			return err
		}
		u.Password = pwd
	}

	return nil
}

// CreateToken returns a jwt.Token with claims from the User.
func (u *User) CreateToken(duration time.Duration) (jwt.Token, error) {
	return jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(duration)).
		Claim("account_id", u.AccountId).
		Claim("username", u.Username).
		Claim("role", u.Role).
		Build()
}
