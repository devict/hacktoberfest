package main

import (
	"log"
	"net/http"

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
		Projects map[string]bool
		User     goth.User
		New      bool
	}{
		Orgs:     orgs,
		Projects: projects,
		User:     u,
		New:      n,
	}

	v.HTML(w, http.StatusOK, "profile", info)
}
