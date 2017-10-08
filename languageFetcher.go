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
// by using a channel of functions to manage the
// fetchedRepos cache.
type languageFetcher struct {
	fetchedRepos map[string][]string
	wg           sync.WaitGroup
	action       chan func()
	stop         chan struct{}
	errors       chan error
}

// newLanguageFetcher makes all the things.
func newLanguageFetcher() *languageFetcher {
	return &languageFetcher{
		fetchedRepos: make(map[string][]string),
		action:       make(chan func()),
		stop:         make(chan struct{}),
		errors:       make(chan error),
	}
}

// repoLanguages gets top three repo languages.
func (lf *languageFetcher) repoLanguages(ctx context.Context, repoURL, token string) {
	defer lf.wg.Done()

	// Check if repo is in cache.
	isCached := make(chan bool)
	lf.action <- func() {
		if langs := lf.fetchedRepos[repoURL]; langs != nil {
			isCached <- true
			return
		}
		isCached <- false
	}

	// If cached don't get languages for repo again
	if <-isCached {
		return
	}

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

	// Cache repo languages.
	actionDone := make(chan bool, 1)
	lf.action <- func() {
		lf.fetchedRepos[repoURL] = langs
		actionDone <- true
	}
	<-actionDone
}

// start starts a loop that runs action functions in the
// order recieved. Also, it has a stop chan just in case.
func (lf *languageFetcher) start() {
	for {
		select {
		case f := <-lf.action:
			f()
		case <-lf.stop:
			return
		}
	}
}

// appendLanguages to deduped issues.
func (lf *languageFetcher) appendLanguages(issues []Issue) (out []Issue) {
	for _, issue := range issues {
		issue.Languages = append(issue.Languages, lf.fetchedRepos[issue.RepoURL]...)
		out = append(out, issue)
	}
	return
}

// wait checks for errors or waits for language workers to finish.
func (lf *languageFetcher) wait(cancel context.CancelFunc) error {
	wgDone := make(chan struct{})
	for {
		select {
		// One of the language workers failed so cancel the others and pass the error up
		case err := <-lf.errors:
			cancel()
			return err
		case <-wgDone:
			return nil
		}
	}
}

// workerWait closes wgDone because wg.Wait() doesn't return a channel.
func workerWait(lf *languageFetcher, wgDone chan struct{}) {
	lf.wg.Wait()
	close(wgDone)
}
