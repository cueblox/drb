package cmd

import (
	"github.com/devrel-blox/drb/blox"
)

type Data struct {
	Profiles []blox.Profile `json:"profiles"`
	Articles []blox.Article `json:"articles"`
}

func NewData() Data {
	var profiles []blox.Profile
	var articles []blox.Article
	return Data{
		Profiles: profiles,
		Articles: articles,
	}

}
