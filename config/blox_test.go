package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
)

var drb_blox_v1 = Schema{
	Namespace: "schemas.devrel-blox.com",
	Name:      "blox",
	Sources: []string{
		"https://schemas.devrel-blox.com/v1",
		"https://altschemas.devrel-blox.com/v1",
	},
	Models: []*Model{
		&drb_profile_v1,
	},
}

var drb_profile_v1 = Model{
	ID:         "profile",
	Name:       "profile",
	Folder:     "profiles",
	ForeignKey: "profile_id",
	Cue: `{
		first_name: string
		last_name:  string
		age?:       int
		company?:   string
		title?:     string
		social_accounts?: [#TwitterAccount | #GitHubAccount | #MiscellaneousAccount]
	
		#TwitterAccount: {
			network:  "twitter"
			username: string
			url:      string | *"https://twitter.com/\(username)"
		}
	
		#GitHubAccount: {
			network:  "github"
			username: string
			url:      string | *"https://github.com/\(username)"
		}
	
		#MiscellaneousAccount: {
			network: string
			url:     string
		}
	}`,
}
var blox_conf = BloxConfig{
	Base:             "testdata",
	Source:           "source",
	Destination:      "out",
	Templates:        "templates",
	Static:           "static",
	Schemas:          "schemas",
	DefaultExtension: ".md",
	SchemaList:       []*Schema{},
}

func TestSchema(t *testing.T) {

	bb, err := yaml.Marshal(drb_blox_v1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(bb))
	jbb, err := json.Marshal(drb_blox_v1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println("-------")
	fmt.Println(string(jbb))
	fmt.Println("-------")
	err = drb_blox_v1.Save()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = blox_conf.AddSchema(&drb_blox_v1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestConfigLoad(t *testing.T) {
	bb, err := yaml.Marshal(blox_conf)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(bb))
}

func TestModel(t *testing.T) {

	bb, err := yaml.Marshal(drb_profile_v1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(bb))
	jbb, err := json.Marshal(drb_profile_v1)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(jbb))

}
