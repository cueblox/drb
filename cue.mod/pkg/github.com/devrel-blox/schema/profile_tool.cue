package profile_tool

import (
	"strings"
	"tool/cli"
	"tool/exec"
	"tool/file"
)

command: {
	profile: {
		profile_dir: cli.Ask & {
			prompt:   "Which directory for profiles?"
			response: string
		}
		find_profiles: file.Glob & {
			glob: profile_dir.response + "/*.yaml"
		}
		print_profiles: cli.Print & {
			text: "Found the following profiles: " + strings.Join(find_profiles.files, ",")
		}
		for i, profile in find_profiles.files {
			"validate-\(i)": exec.Run & {
				cmd:    ["cue", "vet", profile, "profile/profile.cue", "-d", "#Profile"]
				stdout: string
			}
		}
	}
}
