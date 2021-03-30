package models

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Company string `json:"company"`
	Title   string `json:"title"`

	SocialAccounts []Social `json:"social_accounts"`
}

type Social struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	ProfileURL string `json:"profile_url"`
}
