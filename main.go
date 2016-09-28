package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	ghauth "github.com/markbates/goth/providers/github"
	"github.com/unrolled/render"
)

func main() {
	goth.UseProviders(
		ghauth.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK")+"/auth/github/callback"),
	)

	r := pat.New()
	v := render.New(render.Options{
		Layout: "layout",
	})
	gh := github.NewClient(nil)

	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// TODO save user to db
		v.HTML(w, http.StatusOK, "user", user)
	})
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)

	r.Post("/check", func(w http.ResponseWriter, r *http.Request) {
		q := fmt.Sprintf(
			"author:%s+type:pr+created:2016-09-30T00:00:00-12:00..2016-10-31T23:59:59-12:00",
			"jcbwlkr",
		)

		result, resp, err := gh.Search.Issues(q, nil)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("result", result)
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("resp", resp)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		v.HTML(w, http.StatusOK, "home", nil)
	})

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	fmt.Println("Listing on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
