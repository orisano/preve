package preve

import "github.com/google/go-github/github"

func MustGithubClient(baseURL string) *github.Client {
	if baseURL == "" {
		return github.NewClient(nil)
	}
	gh, err := github.NewEnterpriseClient(baseURL, baseURL, nil)
	if err != nil {
		panic(err)
	}
	return gh
}
