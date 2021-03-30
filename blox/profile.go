package blox

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type Profile struct {
	baseModel
	FirstName string `json:"first_name" yaml:"first_name"`
	LastName  string `json:"last_name"  yaml:"last_name"`

	Company string `json:"company" yaml:"company"`
	Title   string `json:"title" yaml:"title"`

	SocialAccounts []SocialAccount `json:"social_accounts" yaml:"social_accounts"`
}

func (p *Profile) Validate() error {
	var err error

	// I'm sure there's a better way to do this, but this is fine for now; right?
	if "" == p.FirstName {
		err = multierror.Append(err, errors.New("Missing first_name"))
	}

	if "" == p.LastName {
		err = multierror.Append(err, errors.New("Missing last_name"))
	}

	if "" == p.Company {
		err = multierror.Append(err, errors.New("Missing company"))
	}

	if "" == p.Title {
		err = multierror.Append(err, errors.New("Missing title"))
	}

	return err
}

type SocialAccount struct {
	Value interface{}
}

func (s *SocialAccount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var typ struct {
		Type string `json:"network"`
	}

	if err := unmarshal(&typ); err != nil {
		return err
	}

	switch strings.ToLower(typ.Type) {
	case "twitter":
		s.Value = new(TwitterAccount)
	case "B":
		s.Value = new(GitHubAccount)
	default:
		s.Value = new(MiscellaneousAccount)
	}

	return unmarshal(s.Value)
}

func (s *SocialAccount) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"network"`
	}

	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	switch strings.ToLower(typ.Type) {
	case "twitter":
		s.Value = new(TwitterAccount)
	case "B":
		s.Value = new(GitHubAccount)
	default:
		s.Value = new(MiscellaneousAccount)
	}

	return json.Unmarshal(data, s.Value)

}

type TwitterAccount struct {
	Username string `json:"username"`
}

type GitHubAccount struct {
	Username string `json:"username"`
}

type MiscellaneousAccount struct {
	Url string `json:"url"`
}
