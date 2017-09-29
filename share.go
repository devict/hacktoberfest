package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func getShare(w http.ResponseWriter, r *http.Request) {
	u, _, ok := findUser(r)
	if !ok {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	var share bool
	if err := db.QueryRow("SELECT share_info FROM users WHERE id = $1", u.UserID).Scan(&share); err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(share); err != nil {
		log.Println(err)
	}
}

func updateShare(w http.ResponseWriter, r *http.Request) {
	u, _, ok := findUser(r)
	if !ok {
		http.Error(w, "you are not logged in", http.StatusUnauthorized)
		return
	}

	var share bool
	if err := json.NewDecoder(r.Body).Decode(&share); err != nil {
		log.Println(err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"UPDATE users SET share_info = $1, updated_at = $2 WHERE id = $3",
		share,
		time.Now(),
		u.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
