package models

import (
	"github.com/knockbox/authentication/pkg/payloads"
)

// UserDetails defines a User(s) details in our database.
type UserDetails struct {
	Id             uint   `db:"id"`
	UserId         uint   `db:"user_id"`
	ProfilePicture string `db:"profile_picture"`
	FullName       string `db:"full_name"`
	GithubURL      string `db:"github_url"`
	TwitterURL     string `db:"twitter_url"`
	WebsiteURL     string `db:"website_url"`
	Verified       bool   `db:"verified"`
}

// ApplyUpdate updates the user with the values from the payloads.UserUpdate.
func (u *UserDetails) ApplyUpdate(payload *payloads.UserDetailsUpdate) {
	if payload.ProfilePicture != nil {
		u.ProfilePicture = *payload.ProfilePicture
	}

	if payload.FullName != nil {
		u.FullName = *payload.FullName
	}

	if payload.GithubURL != nil {
		u.GithubURL = *payload.GithubURL
	}

	if payload.TwitterURL != nil {
		u.TwitterURL = *payload.TwitterURL
	}

	if payload.WebsiteURL != nil {
		u.WebsiteURL = *payload.WebsiteURL
	}
}
