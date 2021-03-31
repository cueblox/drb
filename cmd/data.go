package cmd

import (
	"github.com/devrel-blox/drb/blox/article"
	"github.com/devrel-blox/drb/blox/profile"
)

type Data struct {
	Profiles []profile.Profile `json:"profiles"`
	Articles []article.Article `json:"articles"`
}

func NewData() Data {
	var profiles []profile.Profile
	var articles []article.Article
	return Data{
		Profiles: profiles,
		Articles: articles,
	}
}
