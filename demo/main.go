package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // default port
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Println("Serving on :" + port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}