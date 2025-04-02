package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func analyzeDomain(db *Database, input string) {
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
		db.LogEvent(input, "analysis_failed", map[string]string{
			"error": "Invalid domain: " + input,
		})
		return
	}

	log.Printf("Extracted domain for analysis: %s", domain)

	// Log analysis started
	if err := db.LogEvent(domain, "analysis_started", nil); err != nil {
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
		db.LogEvent(domain, "analysis_failed", map[string]string{
			"error": "Failed to make HTTP request: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Check if response is 200
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 status code received for domain %s: %s", domain, resp.Status)
		db.LogEvent(domain, "analysis_failed", map[string]string{
			"error": "HTTP status code not 200: " + resp.Status,
		})
		return
	}

	log.Printf("Received 200 response for domain %s, reading body", domain)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body for domain %s: %v", domain, err)
		db.LogEvent(domain, "analysis_failed", map[string]string{
			"error": "Failed to read response body: " + err.Error(),
		})
		return
	}

	// Search for Shopify indicators
	bodyStr := string(body)
	if strings.Contains(bodyStr, "myshopify") {
		log.Printf("Found 'myshopify' indicator for domain %s", domain)
		db.LogEvent(domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'myshopify' in page content",
		})
		return
	}

	if strings.Contains(bodyStr, "cdn.shopify.com") {
		log.Printf("Found 'cdn.shopify.com' indicator for domain %s", domain)
		db.LogEvent(domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'cdn.shopify.com' in page content",
		})
		return
	}

	// No Shopify indicators found in main page, try checking the checkout page
	log.Printf("No Shopify indicators found in main page for domain %s, checking checkout page", domain)
	
	// Construct checkout URL
	checkoutURL := "https://checkout." + domain + "/checkout/cn"
	log.Printf("Making HTTP request to checkout URL: %s", checkoutURL)
	
	// Create a new client with shorter timeout for checkout request
	checkoutClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 2 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	
	// Make the request to checkout URL
	checkoutResp, err := checkoutClient.Get(checkoutURL)
	if err != nil {
		log.Printf("Checkout HTTP request failed for domain %s: %v", domain, err)
		db.LogEvent(domain, "analysis_failed", map[string]string{
			"error": "No Shopify indicators found in page content and checkout check failed",
		})
		return
	}
	defer checkoutResp.Body.Close()
	
	// Check for Shopify-specific header
	shopifyID := checkoutResp.Header.Get("x-shopid")
	if shopifyID != "" {
		log.Printf("Found 'x-shopid' header in checkout response for domain %s: %s", domain, shopifyID)
		db.LogEvent(domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'x-shopid' header in checkout page response",
			"shopify_id": shopifyID,
		})
		return
	}
	
	// Check for other Shopify indicators in checkout response
	if strings.Contains(checkoutResp.Header.Get("Server"), "Shopify") {
		log.Printf("Found 'Shopify' in Server header for domain %s", domain)
		db.LogEvent(domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'Shopify' in Server header of checkout page",
		})
		return
	}
	
	// No Shopify indicators found in main page or checkout
	log.Printf("No Shopify indicators found for domain %s in main page or checkout", domain)
	db.LogEvent(domain, "analysis_failed", map[string]string{
		"error": "No Shopify indicators found in page content or checkout page",
	})
}
