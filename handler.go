package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get the URL from the form
		inputURL := r.FormValue("url")
		if inputURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		// Ensure URL has a scheme
		if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
			inputURL = "https://" + inputURL
		}

		// Parse the URL
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		// Extract domain
		domain := parsedURL.Hostname()
		if domain == "" {
			http.Error(w, "Invalid domain", http.StatusBadRequest)
			return
		}

		// Log the form submission
		log.Printf("Form submitted for domain: %s", domain)

		// Redirect to the domain-specific path
		http.Redirect(w, r, "/"+domain, http.StatusSeeOther)
		return
	}

	// For GET requests, show the landing page
	// TODO: Implement landing page HTML
} 