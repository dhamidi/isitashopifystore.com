package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func analyzeDomain(input string) {
	log.Printf("Starting domain analysis for input: %s", input)

	// Extract domain from input
	domain := input
	if strings.Contains(input, "/") {
		// If input contains slashes, try to parse as URL
		parsedURL, err := url.Parse(input)
		if err == nil && parsedURL.Hostname() != "" {
			domain = parsedURL.Hostname()
		}
	}

	if domain == "" {
		log.Printf("Failed to extract valid domain from input: %s", input)
		logEvent(db, input, "analysis_failed", map[string]string{
			"error": "Invalid domain: " + input,
		})
		return
	}

	log.Printf("Extracted domain for analysis: %s", domain)

	// Log analysis started
	if err := logEvent(db, domain, "analysis_started", nil); err != nil {
		log.Printf("Error logging analysis start for domain %s: %v", domain, err)
		return
	}

	// Create URL for analysis
	analysisURL := "https://" + domain
	log.Printf("Making HTTP request to: %s", analysisURL)

	// Create HTTP client that follows redirects (up to 3 times)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				log.Printf("Max redirects reached for domain %s", domain)
				return http.ErrUseLastResponse
			}
			log.Printf("Following redirect %d for domain %s: %s", len(via), domain, req.URL)
			return nil
		},
	}

	// Make the request
	resp, err := client.Get(analysisURL)
	if err != nil {
		log.Printf("HTTP request failed for domain %s: %v", domain, err)
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "Failed to make HTTP request: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Check if response is 200
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 status code received for domain %s: %s", domain, resp.Status)
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "HTTP status code not 200: " + resp.Status,
		})
		return
	}

	log.Printf("Received 200 response for domain %s, reading body", domain)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body for domain %s: %v", domain, err)
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "Failed to read response body: " + err.Error(),
		})
		return
	}

	// Search for Shopify indicators
	bodyStr := string(body)
	if strings.Contains(bodyStr, "myshopify") {
		log.Printf("Found 'myshopify' indicator for domain %s", domain)
		logEvent(db, domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'myshopify' in page content",
		})
		return
	}

	if strings.Contains(bodyStr, "cdn.shopify.com") {
		log.Printf("Found 'cdn.shopify.com' indicator for domain %s", domain)
		logEvent(db, domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'cdn.shopify.com' in page content",
		})
		return
	}

	// No Shopify indicators found
	log.Printf("No Shopify indicators found for domain %s", domain)
	logEvent(db, domain, "analysis_failed", map[string]string{
		"error": "No Shopify indicators found in page content",
	})
} 