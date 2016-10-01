package main

import (
	"fmt"
	"regexp"
)

// Repo is a repository on Github. Owner can be either an organization or user.
type Repo struct {
	Owner string
	Name  string
}

var reRepo = regexp.MustCompile("https://api.github.com/repos/([^/]+)/([^/]+)")

func repoFromURL(url string) (Repo, error) {
	if !reRepo.MatchString(url) {
		return Repo{}, fmt.Errorf("url %q did not match regexp", url)
	}

	matches := reRepo.FindStringSubmatch(url)
	return Repo{
		Name:  matches[2],
		Owner: matches[1],
	}, nil
}
