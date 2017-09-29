package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"
	
	"github.com/pkg/errors"
)

type Issue struct {
	Title string
	Date  time.Time
	URL   string
	Repo  Repo
}

func issues(w http.ResponseWriter, r *http.Request) {
	u, _, ok := findUser(r)
	if !ok {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	issues, err := fetchIssues(u.AccessToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(issues); err != nil {
		log.Println(err)
	}
	
	log.Println(issues)
}

func fetchIssues(token string) ([]Issue, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/search/issues", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not build request")
	}
	
	q := "is:open is:issue label:hacktoberfest"
	
	// add specific organizations to the search
	for k, v := range orgs {
		if v == true {
			q += fmt.Sprintf(" org:%s", k)
		}
	}
	
	// add specific projects to the search
	for k, v := range projects {
		if v == true {
			q += fmt.Sprintf(" repo:%s", k)
		}
	}

	vals := req.URL.Query()
	vals.Add("q", q)
	req.URL.RawQuery = vals.Encode()

	// Use their access token so it counts against their rate limit
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

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
			Title     string    `json:"title"`
			CreatedAt time.Time `json:"created_at"`
			URL       string    `json:"url"`
			RepoURL   string    `json:"repository_url"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "could not decode json")
	}

	issues := []Issue{}
	for _, item := range data.Items {
		issue := Issue{
			Title: item.Title,
			Date:  item.CreatedAt,
			URL:   item.URL,
		}

		issue.Repo, err = repoFromURL(item.RepoURL)
		if err != nil {
			return nil, errors.Wrapf(err, "could not identify repo from %s", item.RepoURL)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}