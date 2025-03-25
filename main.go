package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set up HTTP routes with pattern matching
	mux := http.NewServeMux()
	mux.HandleFunc("/status/", statusHandler)
	mux.HandleFunc("/{domain}", resultPageHandler)
	mux.HandleFunc("/", landingPageHandler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 