package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/pkg/errors"
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

	if err := saveUser(user); err != nil {
		log.Println(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}

func saveUser(u goth.User) error {
	var found int
	err := db.QueryRow("SELECT id FROM users WHERE id = $1", u.UserID).Scan(&found)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "could not check for existing user")
	}

	if err == sql.ErrNoRows {
		_, err := db.Exec(
			"INSERT INTO users (id, username, name, email, avatar, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			u.UserID,
			u.NickName,
			u.Name,
			u.Email,
			u.AvatarURL,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return errors.Wrap(err, "could not insert user")
		}
		return nil
	}

	_, err = db.Exec(
		"UPDATE users SET username = $1, name = $2, email = $3, avatar = $4, updated_at = $5 WHERE id = $6",
		u.NickName,
		u.Name,
		u.Email,
		u.AvatarURL,
		time.Now(),
		u.UserID,
	)
	if err != nil {
		return errors.Wrap(err, "could not update user")
	}

	return nil
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
