package blox

import (
	"testing"

	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
)

func TestProfileDecoding(t *testing.T) {
	expected := Profile{
		FirstName: "Test",
		LastName:  "Name",
		Company:   "MyOrg",
		Title:     "DevRel",
		SocialAccounts: []SocialAccount{
			{
				Network: "linkedin",
				Url:     "https://www.linkedin.com/random",
			},
			{
				Network:  "github",
				Username: "bketelsen",
			},
			{
				Network:  "twitter",
				Username: "rawkode",
			},
		},
	}

	data := `
first_name: Test
last_name: Name
company: MyOrg
title: DevRel
social_accounts:
- network: twitter
  username: rawkode
- network: github
  username: bketelsen
- network: linkedin
  url: https://www.linkedin.com/random
`
	var profile Profile

	err := yaml.Unmarshal([]byte(data), &profile)

	if nil != err {
		t.Error("Failed to unmarshall YAML")
	}

	if diff := deep.Equal(profile, expected); diff != nil {
		t.Error(diff)
	}
}
