package blox

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"cuelang.org/go/cuego"
	"github.com/devrel-blox/drb/cueutils"
	"github.com/goccy/go-yaml"
)

type Profile struct {
	BaseModel `json:",omitempty"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Company string `json:"company"`
	Title   string `json:"title"`

	SocialAccounts []SocialAccount `json:"social_accounts,omitempty"`
}

type SocialAccount struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	Url      string `json:"url,omitempty"`
}

//go:embed profile.cue
var ProfileCue string

func ProfileFromYAML(path string) (Profile, error) {
	err := cuego.Constrain(&Profile{}, ProfileCue)
	if err != nil {
		return Profile{}, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Profile{}, cueutils.UsefulError(err)
	}

	var profile Profile

	err = yaml.Unmarshal(bytes, &profile)
	if err != nil {
		return Profile{}, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&profile)
	if err != nil {
		return Profile{}, cueutils.UsefulError(err)
	}

	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	profile.ID = slug
	fmt.Printf("Profile '%s' validated successfully\n", profile.ID)

	return profile, nil
}
