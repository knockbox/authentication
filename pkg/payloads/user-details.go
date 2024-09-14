package payloads

type UserDetailsUpdate struct {
	ProfilePicture *string `json:"profile_picture,omitempty" validate:"omitempty,http_url"`
	FullName       *string `json:"full_name,omitempty" validate:"omitempty,gte=0,lte=64"`
	GithubURL      *string `json:"github_url,omitempty" validate:"omitempty,http_url"`
	TwitterURL     *string `json:"twitter_url,omitempty" validate:"omitempty,http_url"`
	WebsiteURL     *string `json:"website_url,omitempty" validate:"omitempty,http_url"`
}
