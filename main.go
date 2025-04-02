package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting isitashopifystore.com server")

	// Initialize database
	db, err := NewDatabase("isitashopifystore.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create events table if it doesn't exist
	if err := db.CreateEventsTable(); err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}

	// Create handler with database
	h := &Handler{db: db}

	// Set up HTTP routes with pattern matching
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.png", h.faviconHandler)
	mux.HandleFunc("/status/", h.statusHandler)
	mux.HandleFunc("/{domain}", h.resultPageHandler)
	mux.HandleFunc("/", h.landingPageHandler)

	log.Println("HTTP routes configured:")
	log.Println("  - GET  /              -> Landing page")
	log.Println("  - POST /              -> Process URL submission")
	log.Println("  - GET  /{domain}      -> Show analysis result")
	log.Println("  - GET  /status/{domain} -> Check analysis status")
	log.Println("  - GET  /favicon.png    -> Serve favicon")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 
