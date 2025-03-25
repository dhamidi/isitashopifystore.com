package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting isitashopifystore.com server")

	// Initialize database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create events table if it doesn't exist
	if err := createEventsTable(db); err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}

	// Set up HTTP routes with pattern matching
	mux := http.NewServeMux()
	mux.HandleFunc("/status/", statusHandler)
	mux.HandleFunc("/{domain}", resultPageHandler)
	mux.HandleFunc("/", landingPageHandler)

	log.Println("HTTP routes configured:")
	log.Println("  - GET  /              -> Landing page")
	log.Println("  - POST /              -> Process URL submission")
	log.Println("  - GET  /{domain}      -> Show analysis result")
	log.Println("  - GET  /status/{domain} -> Check analysis status")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 