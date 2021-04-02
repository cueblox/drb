package cmd

import (
	"github.com/devrel-blox/drb/blox"
)

type Data struct {
	Profiles   []map[string]interface{} `json:"profiles"`
	Articles   []blox.Article           `json:"articles"`
	Categories []blox.Category          `json:"categories"`
}

func NewData() Data {
	var profiles []map[string]interface{}
	var articles []blox.Article
	var categories []blox.Category
	return Data{
		Profiles:   profiles,
		Articles:   articles,
		Categories: categories,
	}

}
