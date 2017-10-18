package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

// languageFetcher manages fetching repo languages
// by using a mutex to handle accress to the
// fetchedRepos cache.
type languageFetcher struct {
	wg     sync.WaitGroup
	errors chan error

	mu           sync.Mutex
	fetchedRepos map[string][]string
}

// newLanguageFetcher makes all the things.
func newLanguageFetcher() *languageFetcher {
	return &languageFetcher{
		fetchedRepos: make(map[string][]string),
		errors:       make(chan error),
	}
}

// repoLanguages gets top three repo languages.
func (lf *languageFetcher) repoLanguages(ctx context.Context, repoURL, token string) {
	defer lf.wg.Done()

	lf.mu.Lock()
	if langs := lf.fetchedRepos[repoURL]; langs != nil {
		return
	}
	lf.mu.Unlock()

	// If not cached, get languages from repo.
	req, err := http.NewRequest("GET", fmt.Sprintf(repoURL+"/languages"), nil)
	if err != nil {
		lf.errors <- errors.Wrap(err, "could not build request")
	}

	// Tell the request to use our context so we can cancel it in-flight if needed
	req = req.WithContext(ctx)

	// Use their access token so it counts against their rate limit
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		lf.errors <- errors.Wrap(err, "could not execute request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		lf.errors <- errors.Wrapf(err, "status was %d, not 200", resp.StatusCode)
		return
	}
	data := make(map[string]int)
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		lf.errors <- errors.Wrap(err, "could not decode json")
		return
	}

	// Get top three languages
	langs := top(3, data)

	lf.mu.Lock()
	lf.fetchedRepos[repoURL] = langs
	lf.mu.Unlock()
}

// appendLanguages to deduped issues.
func (lf *languageFetcher) appendLanguages(issues []Issue) (out []Issue) {
	for _, issue := range issues {
		issue.Languages = append(issue.Languages, lf.fetchedRepos[issue.RepoURL]...)
		out = append(out, issue)
	}
	return
}
