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

	// Set up HTTP routes
	http.HandleFunc("/", landingPageHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/", resultPageHandler) // This will handle /[domain] paths

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 