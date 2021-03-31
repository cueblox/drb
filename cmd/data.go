package cmd

import "github.com/devrel-blox/drb/blox/profile"

type Data struct {
	Profiles []profile.Profile `json:"profiles"`
}

func NewData() Data {
	var profiles []profile.Profile
	return Data{
		Profiles: profiles,
	}
}
