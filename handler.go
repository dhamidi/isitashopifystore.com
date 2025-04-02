package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Handler struct {
	db *Database
}

type AnalysisResult struct {
	Status    string `json:"status"`
	Reason    string `json:"reason,omitempty"`
	IsShopify bool   `json:"is_shopify"`
	Domain    string `json:"domain,omitempty"`
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get the URL from the form
		inputURL := strings.TrimSpace(r.FormValue("url"))
		if inputURL == "" {
			log.Printf("Empty URL submitted")
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		log.Printf("Processing URL submission: %s", inputURL)

		// Basic input sanitization
		inputURL = strings.ToLower(inputURL)
		inputURL = strings.TrimPrefix(inputURL, "www.")

		// If only a domain name is submitted, treat it as https
		if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
			// Check if it's a valid domain name
			if !isValidDomain(inputURL) {
				log.Printf("Invalid domain name submitted: %s", inputURL)
				http.Error(w, "Invalid domain name", http.StatusBadRequest)
				return
			}
			inputURL = "https://" + inputURL
		}

		// Parse and validate the URL
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			log.Printf("Invalid URL format: %s, error: %v", inputURL, err)
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Additional URL validation
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			log.Printf("Invalid URL scheme: %s", parsedURL.Scheme)
			http.Error(w, "Only http and https URLs are allowed", http.StatusBadRequest)
			return
		}

		// Extract and validate domain
		domain := parsedURL.Hostname()
		if domain == "" || !isValidDomain(domain) {
			log.Printf("Invalid domain extracted from URL: %s", inputURL)
			http.Error(w, "Invalid domain name", http.StatusBadRequest)
			return
		}

		// Log the form submission
		log.Printf("Form submitted successfully for domain: %s", domain)

		// Redirect to the domain-specific path
		log.Printf("Redirecting to /%s", domain)
		http.Redirect(w, r, "/"+domain, http.StatusSeeOther)
		return
	}

	// For GET requests, show the landing page
	t := template.Must(template.ParseFiles("html/landing_page.html"))
	t.Execute(w, nil)
}

func resultPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract domain from path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		log.Printf("Empty path received")
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	// Extract domain from path
	domain := path
	if strings.Contains(path, "/") {
		// If path contains slashes, try to parse as URL
		parsedURL, err := url.Parse(path)
		if err == nil && parsedURL.Hostname() != "" {
			domain = parsedURL.Hostname()
		}
	}

	if domain == "" {
		log.Printf("Invalid domain extracted from path: %s", path)
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	log.Printf("Processing result page request for domain: %s", domain)

	// Check if analysis exists
	var result AnalysisResult
	var err error
	result.Status, result.Reason, err = db.GetLatestAnalysisResult(domain)

	if err == sql.ErrNoRows {
		log.Printf("No analysis found for domain: %s, starting new analysis", domain)
		// No analysis exists, trigger background analysis
		go analyzeDomain(domain)

		// Show polling page
		log.Printf("Rendering polling page for domain: %s", domain)
		t := template.Must(template.ParseFiles("html/polling_page.html"))
		t.Execute(w, struct{ Domain string }{domain})
		return
	}

	if err != nil {
		log.Printf("Error checking analysis status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse the payload to get the reason
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(result.Reason), &payload); err != nil {
		log.Printf("Error parsing payload: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Determine if it's a Shopify store
	result.IsShopify = result.Status == "analysis_succeeded"
	result.Domain = domain

	// Render the result page
	t := template.Must(template.ParseFiles("html/result_page.html"))
	t.Execute(w, result)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Extract domain from path
	path := strings.TrimPrefix(r.URL.Path, "/status/")
	if path == "" {
		log.Printf("Empty path received in status handler")
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	// Extract domain from path
	domain := path
	if strings.Contains(path, "/") {
		// If path contains slashes, try to parse as URL
		parsedURL, err := url.Parse(path)
		if err == nil && parsedURL.Hostname() != "" {
			domain = parsedURL.Hostname()
		}
	}

	if domain == "" {
		log.Printf("Invalid domain extracted from path in status handler: %s", path)
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	log.Printf("Processing status request for domain: %s", domain)

	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Check if analysis exists
	var result AnalysisResult
	var err error
	result.Status, result.Reason, err = db.GetStatusResult(domain)
	
	if err == sql.ErrNoRows {
		log.Printf("No analysis found for domain in status check: %s", domain)
		// Start a new analysis in background
		go analyzeDomain(domain)
		// Return in_progress status
		json.NewEncoder(w).Encode(AnalysisResult{
			Status: "in_progress",
		})
		return
	}

	if err != nil {
		log.Printf("Error checking analysis status for domain %s: %v", domain, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Map event types to status
	switch result.Status {
	case "analysis_started":
		result.Status = "in_progress"
		log.Printf("Analysis in progress for domain: %s", domain)
	case "analysis_succeeded":
		result.Status = "succeeded"
		result.IsShopify = true
		log.Printf("Analysis succeeded for domain: %s", domain)
		// Parse the reason from the payload
		var payload map[string]string
		if err := json.Unmarshal([]byte(result.Reason), &payload); err == nil {
			if reason, ok := payload["reason"]; ok {
				result.Reason = reason
			}
		}
	case "analysis_failed":
		result.Status = "failed"
		log.Printf("Analysis failed for domain: %s", domain)
		// Parse the error from the payload
		var payload map[string]string
		if err := json.Unmarshal([]byte(result.Reason), &payload); err == nil {
			if errorMsg, ok := payload["error"]; ok {
				result.Reason = errorMsg
			}
		}
	default:
		result.Status = "in_progress"
		log.Printf("Unknown analysis status for domain %s: %s", domain, result.Status)
	}

	// Return the result
	json.NewEncoder(w).Encode(result)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("assets", "favicon-512.png")

	// Read the favicon file
	favicon, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading favicon: %v", err)
		http.Error(w, "Favicon not found", http.StatusNotFound)
		return
	}

	// Set content type and serve
	w.Header().Set("Content-Type", "image/png")
	w.Write(favicon)
}

// isValidDomain checks if a string is a valid domain name
func isValidDomain(domain string) bool {
	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$`)
	return domainRegex.MatchString(domain)
}
