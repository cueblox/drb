{
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
}