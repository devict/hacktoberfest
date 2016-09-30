package main

import "net/http"

func profile(w http.ResponseWriter, r *http.Request) {
	u, ok := findUser(r)
	if !ok {
		http.Redirect(w, r, "/auth/github", http.StatusTemporaryRedirect)
		return
	}

	v.HTML(w, http.StatusOK, "profile", u)
}
