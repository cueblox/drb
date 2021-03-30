package blox

import (
	"testing"

	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
)

func TestValidate(t *testing.T) {
	validProfile := Profile{
		FirstName: "David",
		LastName:  "McKay",
		Company:   "Rawkode Enterprises",
		Title:     "Chief Chief Officer",
	}

	if len(validProfile.Validate()) > 0 {
		t.Error("Valid Profile failed to validate")
	}

	invalidProfile := Profile{}

	if len(invalidProfile.Validate()) == 0 {
		t.Error("Invalid Profile failed to invalidate")
	}
}

func TestProfileDecoding(t *testing.T) {
	expected := Profile{
		FirstName: "Test",
		LastName:  "Name",
		Company:   "MyOrg",
		Title:     "DevRel",
		SocialAccounts: []SocialAccount{
			{Value: MiscellaneousAccount{
				Url: "https://www.linkedin.com/random",
			}},
			{Value: GitHubAccount{
				Username: "bketelsen",
			}},
			{Value: TwitterAccount{
				Username: "rawkode",
			}},
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

	// spew.Println(profile)

	if nil != err {
		t.Error("Failed to unmarshall YAML")
	}

	if diff := deep.Equal(profile, expected); diff != nil {
		t.Error(diff)
	}
}
