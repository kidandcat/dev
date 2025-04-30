package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004" // default port
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Add handler for /hello endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Add handler for /bye endpoint
	http.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Write([]byte("bye " + q))
	})

	log.Println("Serving on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
