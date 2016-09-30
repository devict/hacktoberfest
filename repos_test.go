package main

import "testing"

func TestRepoFromURL(t *testing.T) {
	tests := []struct {
		url  string
		repo Repo
		ok   bool
	}{
		{"https://api.github.com/repos/spf13/hugo", Repo{Owner: "spf13", Name: "hugo"}, true},
		{"https://api.github.com/repos/exercism/xgo", Repo{Owner: "exercism", Name: "xgo"}, true},
		{"https://api.github.com/", Repo{}, false},
	}

	for i, test := range tests {
		got, err := repoFromURL(test.url)
		if test.ok && err != nil {
			t.Errorf("%d: error should be nil, got %v", i, err)
		} else if !test.ok && err == nil {
			t.Errorf("%d: error should not be nil, but it was", i)
		} else if got != test.repo {
			t.Errorf("%d: got != want:\n%+v\n%+v", i, got, test.repo)
		}
	}
}
