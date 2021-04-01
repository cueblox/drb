package cmd

import (
	"github.com/devrel-blox/drb/blox"
)

type Data struct {
	Profiles   []blox.Profile  `json:"profiles"`
	Articles   []blox.Article  `json:"articles"`
	Categories []blox.Category `json:"categories"`
}

func NewData() Data {
	var profiles []blox.Profile
	var articles []blox.Article
	var categories []blox.Category
	return Data{
		Profiles:   profiles,
		Articles:   articles,
		Categories: categories,
	}

}
