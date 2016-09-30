package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/markbates/goth"
	"github.com/pkg/errors"
)

type PR struct {
	URL   string
	Repo  Repo
	Valid bool
}

type Repo struct {
	Owner string
	Name  string
}

// Any project under one of these organizations counts
var orgs = map[string]bool{
	"devICT":         true,
	"MakeICT":        true,
	"openwichita":    true,
	"startupwichita": true,
	"wichitalks":     true,
	"Ennovar":        true,
}

// These specific projects also count
var projects = map[string]bool{
	"imacrayon/foodtrucksnear.me": true,
}

func check(w http.ResponseWriter, r *http.Request) {
	u, ok := findUser(r)
	if !ok {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	prs, err := fetchPRs(u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, pr := range prs {
		prs[i].Valid = orgs[pr.Repo.Owner] || projects[pr.Repo.Owner+"/"+pr.Repo.Name]
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prs); err != nil {
		log.Println(err)
	}
}

func fetchPRs(u goth.User) ([]PR, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/search/issues", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not build request")
	}

	q := fmt.Sprintf(
		"author:%s type:pr created:2016-08-30T00:00:00-12:00..2016-10-31T23:59:59-12:00", // TODO fix dates
		u.NickName,
	)
	vals := req.URL.Query()
	vals.Add("q", q)
	req.URL.RawQuery = vals.Encode()

	req.Header.Add("Authorization", "token "+u.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(err, "status was %d, not 200", resp.StatusCode)
	}

	var data struct {
		Items []struct {
			URL     string `json:"url"`
			RepoURL string `json:"repository_url"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "could not decode json")
	}

	var prs []PR
	for _, item := range data.Items {
		pr := PR{URL: item.URL}

		pr.Repo, err = repoFromURL(item.RepoURL)
		if err != nil {
			return nil, errors.Wrapf(err, "could not identify repo from %s", item.RepoURL)
		}

		prs = append(prs, pr)
	}

	return prs, nil
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
