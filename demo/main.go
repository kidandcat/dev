package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004" // default port
	}

	db, err := sql.Open("sqlite3", "file:names.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS names (name TEXT)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))

	// /hello endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// /bye endpoint
	http.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name != "" {
			_, err := db.Exec("INSERT INTO names(name) VALUES(?)", name)
			if err != nil {
				log.Printf("Failed to insert name: %v", err)
			}
		}
		w.Write([]byte("bye " + name))
	})

	// /names endpoint
	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT name FROM names")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to query names"))
			return
		}
		defer rows.Close()
		names := []string{}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				continue
			}
			names = append(names, name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(names)
	})

	log.Println("Serving on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
