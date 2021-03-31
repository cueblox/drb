package cmd

import "github.com/devrel-blox/drb/blox"

type Data struct {
	Profiles []blox.Profile `json:"profiles"`
}

func NewData() Data {
	var profiles []blox.Profile
	return Data{
		Profiles: profiles,
	}
}
