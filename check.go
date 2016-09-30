package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/goth"
	"github.com/pkg/errors"
)

type Repo struct {
	Owner string
	Name  string
}

var repos map[string]Repo

func init() {
	repos = make(map[string]Repo)
}

type CheckResult struct {
	Valid   []Repo
	Invalid []Repo
}

// Any project under one of these organizations counts
var orgs = map[string]bool{
	"devICT":         true,
	"MakeICT":        true,
	"openwichita":    true,
	"startupwichita": true,
	"wichitalks":     true,
}

// These specific projects also count
var projects = map[string]bool{
	"imacrayon/wichitawesome": true,
}

func check(w http.ResponseWriter, r *http.Request) {
	u, ok := findUser(r)
	if !ok {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	found, err := prRepos(u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := CheckResult{
		Valid:   []Repo{},
		Invalid: []Repo{},
	}

	for _, repo := range found {
		if orgs[repo.Owner] || projects[repo.Owner+"/"+repo.Name] {
			res.Valid = append(res.Valid, repo)
		} else {
			res.Invalid = append(res.Invalid, repo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}

func prRepos(u goth.User) ([]Repo, error) {
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

	var prs struct {
		Items []struct {
			RepoURL string `json:"repository_url"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, errors.Wrap(err, "could not decode json")
	}

	var found []Repo
	for _, pr := range prs.Items {
		r, err := fetchRepo(pr.RepoURL, u)
		if err != nil {
			return nil, errors.Wrapf(err, "could not fetch repo %s", pr.RepoURL)
		}

		found = append(found, r)
	}

	return found, nil
}

func fetchRepo(url string, u goth.User) (Repo, error) {
	if r, ok := repos[url]; ok {
		return r, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Repo{}, errors.Wrap(err, "could not build request")
	}
	req.Header.Add("Authorization", "token "+u.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Repo{}, errors.Wrap(err, "could not execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Repo{}, errors.Wrapf(err, "status was %d, not 200", resp.StatusCode)
	}

	var info struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return Repo{}, errors.Wrap(err, "could not decode json")
	}

	r := Repo{
		Owner: info.Owner.Login,
		Name:  info.Name,
	}

	repos[url] = r
	return r, nil
}
