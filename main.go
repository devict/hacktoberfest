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

type Organization struct {
	Title     string
	ImagePath string
	Visible   bool
}

// Any project under one of these organizations counts
var orgs = map[string]Organization{
	"devict": {
		Title:     "devICT",
		ImagePath: "public/images/logos/devict.svg",
		Visible:   true,
	},
	"grooverlabs": {
		Title:     "Groover Labs",
		ImagePath: "public/images/logos/grooverlabs.png",
		Visible:   true,
	},
	"makeict": {
		Title:     "makeICT",
		ImagePath: "public/images/logos/makeict.svg",
		Visible:   true,
	},
	"lake-afton-public-observatory": {
		Title:     "Lake Afton Public Observatory",
		ImagePath: "public/images/logos/lake-afton.png",
		Visible:   true,
	},
}

type Sponsor struct {
	Name         string
	URL          string
	ImagePath    string
	ImageClasses string
}

var sponsors = []Sponsor{
	{
		Name:         "Moonbase Labs",
		URL:          "https://moonbaselabs.com",
		ImagePath:    "public/images/logos/moonbaselabs.svg",
		ImageClasses: "",
	},
	{
		Name:         "Quilibrium",
		URL:          "https://quilibrium.com/",
		ImagePath:    "public/images/logos/quilibrium.svg",
		ImageClasses: "",
	},
}

type Project struct {
	Title       string
	Description string
	Visible     bool
}

type HacktoberfestConfiguration struct {
	RequiredPullRequestCount    int
	RequiredPullRequestCountEng string
	TimesToRepeat               string
	DisplaySponsors             bool
}

var hacktoberfestOptions = HacktoberfestConfiguration{
	RequiredPullRequestCount:    2,                      // the number of pull requests required to qualify for this year.
	RequiredPullRequestCountEng: "two pull requests",    // the required number written out fully
	TimesToRepeat:               "Repeat 1 more times.", // the value in this sentence should be one less than the required count.
	DisplaySponsors:             false,
}

// These specific projects also count
var projects = map[string]Project{
	"imacrayon/eventsinwichita": {
		Title:       "Events in Wichita",
		Description: "",
		Visible:     false,
	},
	"br0xen/boltbrowser": {
		Title:       "Bolt Browser",
		Description: "A CLI Browser for BoltDB Files.",
		Visible:     false,
	},
	"benblankley/fort-rpg": {
		Title:       "fort-rpg",
		Description: "A text-based Computer Role Playing Game written in Fortran 90.",
		Visible:     true,
	},
	"blunket/camelot": {
		Title:       "Camelot",
		Description: "The 2-player strategy board game, Camelot! (a.k.a. Inside Moves)",
		Visible:     true,
	},
	"kentonh/ProjectNameGenerator": {
		Title:       "Project Name Generator",
		Description: "Really stupid way to give your project a code name.",
		Visible:     true,
	},
	"sethetter/reqq": {
		Title:       "reqq",
		Description: "A rust CLI for making predefined HTTP requests.",
		Visible:     true,
	},
	"sethetter/linktrap": {
		Title:       "linktrap",
		Description: "Text a link to a twilio-number, get back an archived and unpaywalled version.",
		Visible:     true,
	},
	"imacrayon/whatthetofu": {
		Title:       "What the Tofu",
		Description: "Find vegan food in Wichita.",
		Visible:     true,
	},
	"imacrayon/alpine-ajax": {
		Title:       "Alpine AJAX",
		Description: "An Alpine.js plugin for building AJAX-powered frontends.",
		Visible:     true,
	},
	"imacrayon/snowbodyknows": {
		Title:       "Snowbody Knows",
		Description: "A secret santa wishlist builder.",
		Visible:     true,
	},
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
		Orgs     map[string]Organization
		Sponsors []Sponsor
		Projects map[string]Project
		Year     int
		Config   HacktoberfestConfiguration
	}{
		Orgs:     orgs,
		Sponsors: sponsors,
		Projects: projects,
		Year:     time.Now().Year(),
		Config:   hacktoberfestOptions,
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
