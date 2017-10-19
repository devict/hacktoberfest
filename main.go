package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/pat"
	_ "github.com/lib/pq"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
)

// Any project under one of these organizations counts
var orgs = map[string]bool{
	"devict":         true,
	"MakeICT":        true,
	"openwichita":    true,
	"StartupWichita": true,
	"Wichitalks":     true,
	"Ennovar":        true,
}

// These specific projects also count
var projects = map[string]bool{
	"imacrayon/eventsinwichita": true,
	"br0xen/boltbrowser":        true,
	"benblankley/fort-rpg":      true,
	"chrisl8/ArloBot":           true,
	"blunket/image-checker":     true,
}

var v = render.New(render.Options{
	Layout:        "layout",
	IsDevelopment: dev(),
})

var db *sql.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := setupDB(); err != nil {
		log.Println("could not set up db", err)
		os.Exit(1)
	}

	run := "web"
	if len(os.Args) > 1 {
		run = os.Args[1]
	}

	if run == "check" {
		if err := check(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// default: web

	r := pat.New()

	// Register auth handlers. pat requires all routes be registered most
	// specific first so the shorter routes have to be added last
	r.Get("/auth/{provider}/callback", authCallback)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)

	r.Get("/api/issues", issues)
	r.Get("/api/prs", prs)
	r.Get("/api/share", getShare)
	r.Put("/api/share", updateShare)

	r.Get("/profile", profile)

	// Serve static files
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	r.Get("/", home)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	fmt.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, logger(r)))
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Serving %s", r.URL.String())
		h.ServeHTTP(w, r)
		log.Printf("Done serving %s [%v]", r.URL.String(), time.Since(start))
	})
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

func setupDB() error {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return errors.Wrap(err, "could not open db")
	}

	// See if the db is awake. Give it a while to come up with a simple backoff
	for attempts := 1; attempts <= 20; attempts++ {
		if err = db.Ping(); err == nil {
			break
		}
		log.Printf("Waiting for db to come up, attempt %d of 20: %v", attempts, err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		return errors.Wrap(err, "db never came up")
	}

	q := `CREATE TABLE IF NOT EXISTS users (
		id integer,
		username varchar(255),
		name varchar(255),
		email varchar(255),
		avatar varchar(255),
		share_info boolean default false,
		created_at TIMESTAMP WITH TIME ZONE,
		updated_at TIMESTAMP WITH TIME ZONE,
		PRIMARY KEY(id)
	)`
	_, err = db.Exec(q)
	if err != nil {
		return errors.Wrap(err, "could not make user table")
	}

	return nil
}
