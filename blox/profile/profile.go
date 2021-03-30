package profile

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/yaml"
	"github.com/devrel-blox/drb/blox"
)

type Profile struct {
	blox.BaseModel

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Company string `json:"company"`
	Title   string `json:"title"`

	SocialAccounts []SocialAccount `json:"social_accounts"`
}

type SocialAccount struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	Url      string `json:"url"`
}

const CUE = `first_name: "David" | "Brian"
last_name: string
company?: string
title?: string
social_accounts?: [#TwitterAccount | #GitHubAccount | #MiscellaneousAccount]

#TwitterAccount :: {
	network: "twitter"
	username: string
	url: string | *"https://twitter.com/"
}

#GitHubAccount :: {
	network: "github"
	username: string
	url: string | *"https://github.com/"
}

#MiscellaneousAccount :: {
	network: string
	url: string
}
`

func LoadFromYAML(path string) (Profile, error) {
	var cueRuntime cue.Runtime
	profileInstance, err := cueRuntime.Compile("profile", CUE)

	if err != nil {
		return Profile{}, err
	}

	valueInstance, err := yaml.Decode(&cueRuntime, path, nil)
	if err != nil {
		return Profile{}, fmt.Errorf("Parse YAML file error: %w", err)
	}

	merged := cue.Merge(profileInstance, valueInstance)
	err = merged.Value().Validate()
	if err != nil {
		return Profile{}, fmt.Errorf("Validation error: %w", err)
	}

	var profile Profile
	codec := gocodec.New(&cueRuntime, &gocodec.Config{})

	err = codec.Encode(merged.Value(), &profile)

	if err != nil {
		return Profile{}, fmt.Errorf("Encode error: %w", err)
	}

	return profile, nil
}
