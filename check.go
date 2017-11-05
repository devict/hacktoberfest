package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func check() error {
	fmt.Println("Checking results.")

	var users []string
	rows, err := db.Query("SELECT username FROM users ORDER BY username")

	if err != nil {
		return errors.Wrap(err, "could not query list of users")
	}

	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return errors.Wrap(err, "could not scan user")
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "could not iterate over user result")
	}

	prs, err := fetchPRs(users, os.Getenv("PAT"))
	if err != nil {
		return errors.Wrap(err, "could not fetch PRs")
	}

	prsByUser := make(map[string][]PR)
	for _, p := range prs {
		prsByUser[p.User] = append(prsByUser[p.User], p)
	}

	var succesful, unsuccessful int
	fmt.Printf("   %20s %8s %8s\n", "Username", "Valid", "Invalid")
	for i, u := range users {
		prs := prsByUser[u]

		var valid, invalid int
		for _, p := range prs {
			if p.Valid {
				valid++
			} else {
				invalid++
			}
		}
		var result string
		if valid >= 4 {
			succesful++
			result = "✔"
		} else {
			unsuccessful++
			result = "✘"
		}
		fmt.Printf("%02d %20s %8d %8d %s\n", i, u, valid, invalid, result)
	}
	fmt.Printf("Succesful: %d\nUnsuccessful: %d\n", succesful, unsuccessful)

	return nil
}
