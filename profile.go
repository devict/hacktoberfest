package main

import (
	"log"
	"net/http"
	"time"

	"github.com/markbates/goth"
)

func profile(w http.ResponseWriter, r *http.Request) {
	u, n, ok := findUser(r)
	if !ok {
		http.Redirect(w, r, "/auth/github", http.StatusTemporaryRedirect)
		return
	}

	s, _ := sess.Get(r, "session")
	s.Values["new"] = false
	if err := s.Save(r, w); err != nil {
		log.Println(err)
	}

	info := struct {
		Orgs     map[string]bool
		Projects map[string]Project
		User     goth.User
		New      bool
		Config   HacktoberfestConfiguration
		Year     int
	}{
		Orgs:     orgs,
		Projects: projects,
		User:     u,
		New:      n,
		Config:   hacktoberfestOptions,
		Year:     time.Now().Year(),
	}

	v.HTML(w, http.StatusOK, "profile", info)
}
