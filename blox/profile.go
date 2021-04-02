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

type SocialAccount struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	Url      string `json:"url,omitempty"`
}

//go:embed profile.cue
var ProfileCue string

//go:embed pat.md
var ProfileTemplate string

func ProfileFromYAML(path string, schema string) (map[string]interface{}, error) {
	empty := make(map[string]interface{})

	var profile map[string]interface{}

	err := cuego.Constrain(&profile, schema)
	if err != nil {
		return empty, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return empty, cueutils.UsefulError(err)
	}

	err = yaml.Unmarshal(bytes, &profile)
	if err != nil {
		return empty, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&profile)
	if err != nil {
		return empty, cueutils.UsefulError(err)
	}

	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	profile["ID"] = slug
	fmt.Printf("Profile '%s' validated successfully\n", profile["ID"])

	return profile, nil
}
