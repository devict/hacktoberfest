package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/markbates/goth/gothic"
	"github.com/unrolled/render"
)

var v = render.New(render.Options{
	Layout: "layout",
})

func main() {
	r := pat.New()

	// Register auth handlers. pat requires all routes be registered most
	// specific first so the shorter routes have to be added last
	r.Get("/auth/{provider}/callback", authCallback)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)

	r.Get("/profile", profile)
	r.Get("/api/check", check) // TODO change to post

	r.Get("/", home)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	fmt.Println("Listing on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func home(w http.ResponseWriter, r *http.Request) {
	v.HTML(w, http.StatusOK, "home", nil)
}
