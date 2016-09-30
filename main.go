package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/markbates/goth/gothic"
	"github.com/unrolled/render"
)

// Any project under one of these organizations counts
var orgs = map[string]bool{
	"devICT":         true,
	"MakeICT":        true,
	"openwichita":    true,
	"startupwichita": true,
	"wichitalks":     true,
	"Ennovar":        true,
}

// These specific projects also count
var projects = map[string]bool{
	"imacrayon/foodtrucksnear.me": true,
}

var v = render.New(render.Options{
	Layout:        "layout",
	IsDevelopment: dev(),
})

func main() {
	r := pat.New()

	// Register auth handlers. pat requires all routes be registered most
	// specific first so the shorter routes have to be added last
	r.Get("/auth/{provider}/callback", authCallback)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)

	r.Get("/profile", profile)
	r.Get("/api/prs", prs)

	r.Get("/", home)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	fmt.Println("Listing on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Orgs     map[string]bool
		Projects map[string]bool
	}{
		Orgs:     orgs,
		Projects: projects,
	}
	v.HTML(w, http.StatusOK, "home", data)
}

func dev() bool {
	d, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		return false
	}

	return d
}
