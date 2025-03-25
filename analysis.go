package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func analyzeDomain(domain string) {
	// Log analysis started
	if err := logEvent(db, domain, "analysis_started", nil); err != nil {
		log.Printf("Error logging analysis start: %v", err)
		return
	}

	// Ensure domain has https scheme
	if !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	// Parse URL to ensure we only analyze the root path
	parsedURL, err := url.Parse(domain)
	if err != nil {
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "Invalid domain: " + err.Error(),
		})
		return
	}
	parsedURL.Path = "/"

	// Create HTTP client that follows redirects (up to 3 times)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Make the request
	resp, err := client.Get(parsedURL.String())
	if err != nil {
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "Failed to make HTTP request: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Check if response is 200
	if resp.StatusCode != http.StatusOK {
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "HTTP status code not 200: " + resp.Status,
		})
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logEvent(db, domain, "analysis_failed", map[string]string{
			"error": "Failed to read response body: " + err.Error(),
		})
		return
	}

	// Search for Shopify indicators
	bodyStr := string(body)
	if strings.Contains(bodyStr, "myshopify") {
		logEvent(db, domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'myshopify' in page content",
		})
		return
	}

	if strings.Contains(bodyStr, "cdn.shopify.com") {
		logEvent(db, domain, "analysis_succeeded", map[string]string{
			"reason": "Found 'cdn.shopify.com' in page content",
		})
		return
	}

	// No Shopify indicators found
	logEvent(db, domain, "analysis_failed", map[string]string{
		"error": "No Shopify indicators found in page content",
	})
} 