package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

var sess sessions.Store

func init() {
	sess = sessions.NewCookieStore(
		[]byte(os.Getenv("SESSION_SECRET")),
	)

	gothic.Store = sess

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK")+"/auth/github/callback"),
	)
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	s, _ := sess.Get(r, "session")
	s.Values["user"] = user
	if err := s.Save(r, w); err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}

func findUser(r *http.Request) (goth.User, bool) {
	s, err := sess.Get(r, "session")
	if err != nil {
		return goth.User{}, false
	}

	val, ok := s.Values["user"]
	if !ok {
		return goth.User{}, false
	}
	u, ok := val.(goth.User)

	return u, ok
}
