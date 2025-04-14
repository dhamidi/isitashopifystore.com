package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// performSingleAnalysis checks a single URL for Shopify indicators.
// It returns true and a success payload if indicators are found,
// otherwise false and an error payload.
func performSingleAnalysis(analysisURL string) (isShopify bool, payload map[string]string) {
	log.Printf("Performing analysis for URL: %s", analysisURL)
	domain := ""
	parsedURL, err := url.Parse(analysisURL)
	if err == nil && parsedURL.Hostname() != "" {
		domain = parsedURL.Hostname()
	} else {
		log.Printf("Could not parse domain from URL: %s", analysisURL)
		return false, map[string]string{
			"error": fmt.Sprintf("Invalid analysis URL: %s", analysisURL),
		}
	}

	// Create HTTP client that follows redirects (up to 10 times)
	client := &http.Client{
		Timeout: 20 * time.Second, // Add a general timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				log.Printf("Max redirects reached for %s", analysisURL)
				return http.ErrUseLastResponse
			}
			log.Printf("Following redirect %d for %s: %s", len(via), analysisURL, req.URL)
			return nil
		},
	}

	// Make the request to the main page
	resp, err := client.Get(analysisURL)
	if err != nil {
		log.Printf("HTTP request failed for %s: %v", analysisURL, err)
		return false, map[string]string{
			"error": fmt.Sprintf("Failed to make HTTP request to %s: %v", analysisURL, err),
		}
	}
	defer resp.Body.Close()

	// Check if response is 200
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 status code received for %s: %s", analysisURL, resp.Status)
		return false, map[string]string{
			"error": fmt.Sprintf("HTTP status code not 200 for %s: %s", analysisURL, resp.Status),
		}
	}

	log.Printf("Received 200 response for %s, reading body", analysisURL)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body for %s: %v", analysisURL, err)
		return false, map[string]string{
			"error": fmt.Sprintf("Failed to read response body for %s: %v", analysisURL, err),
		}
	}

	// Search for Shopify indicators in the body
	bodyStr := string(body)
	if strings.Contains(bodyStr, "myshopify") {
		log.Printf("Found 'myshopify' indicator for %s", analysisURL)
		return true, map[string]string{
			"reason":    fmt.Sprintf("Found 'myshopify' in page content of %s", analysisURL),
			"indicator": "myshopify",
			"source":    "body",
		}
	}

	if strings.Contains(bodyStr, "cdn.shopify.com") {
		log.Printf("Found 'cdn.shopify.com' indicator for %s", analysisURL)
		return true, map[string]string{
			"reason":    fmt.Sprintf("Found 'cdn.shopify.com' in page content of %s", analysisURL),
			"indicator": "cdn.shopify.com",
			"source":    "body",
		}
	}

	// No Shopify indicators found in main page, try checking the checkout page
	log.Printf("No Shopify indicators found in main page for %s, checking checkout page", analysisURL)

	// Construct checkout URL if possible
	if domain == "" {
		// Cannot construct checkout URL without a valid domain
		return false, map[string]string{
			"error": "No Shopify indicators found in page content, and could not check checkout page (invalid domain)",
		}
	}
	checkoutURL := "https://checkout." + domain + "/checkout/cn"
	log.Printf("Making HTTP request to checkout URL: %s", checkoutURL)

	// Create a new client with shorter timeout for checkout request
	checkoutClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow fewer redirects for checkout probe
			if len(via) >= 2 {
				log.Printf("Max redirects reached for checkout probe %s", checkoutURL)
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Make the request to checkout URL
	checkoutResp, err := checkoutClient.Get(checkoutURL)
	if err != nil {
		log.Printf("Checkout HTTP request failed for %s (derived from %s): %v", checkoutURL, analysisURL, err)
		// This specific error is not conclusive proof it's NOT shopify, just that the checkout probe failed.
		return false, map[string]string{
			"error":  "No Shopify indicators found in main page content.",
			"detail": fmt.Sprintf("Checkout probe failed for %s: %v", checkoutURL, err),
		}
	}
	defer checkoutResp.Body.Close()

	// Check for Shopify-specific header
	shopifyID := checkoutResp.Header.Get("x-shopid")
	if shopifyID != "" {
		log.Printf("Found 'x-shopid' header in checkout response for %s (derived from %s): %s", checkoutURL, analysisURL, shopifyID)
		return true, map[string]string{
			"reason":     fmt.Sprintf("Found 'x-shopid' header in checkout page response (%s)", checkoutURL),
			"indicator":  "x-shopid",
			"source":     "header",
			"shopify_id": shopifyID,
		}
	}

	// Check for other Shopify indicators in checkout response header
	serverHeader := checkoutResp.Header.Get("Server")
	if strings.Contains(serverHeader, "Shopify") {
		log.Printf("Found 'Shopify' in Server header for %s (derived from %s)", checkoutURL, analysisURL)
		return true, map[string]string{
			"reason":    fmt.Sprintf("Found 'Shopify' in Server header of checkout page (%s)", checkoutURL),
			"indicator": "Server: Shopify",
			"source":    "header",
		}
	}

	// No Shopify indicators found anywhere
	log.Printf("No Shopify indicators found for %s in main page or checkout", analysisURL)
	return false, map[string]string{
		"error": "No Shopify indicators found in page content or checkout page.",
	}
}

func analyzeDomain(db *Database, input string) {
	log.Printf("Starting domain analysis for input: %s", input)

	// 1. Extract base domain (remove www. if present)
	var baseDomain string
	potentialURL, err := url.Parse(input)
	if err == nil && potentialURL.Hostname() != "" {
		baseDomain = potentialURL.Hostname()
	} else if strings.Contains(input, ".") && !strings.Contains(input, "/") {
		// Assume it's a domain if it has a dot and no slash
		baseDomain = input
	} else {
		// Try parsing with a default scheme
		potentialURLWithScheme, errScheme := url.Parse("https://" + input)
		if errScheme == nil && potentialURLWithScheme.Hostname() != "" {
			baseDomain = potentialURLWithScheme.Hostname()
		}
	}

	// Clean the base domain (remove www.)
	baseDomain = strings.TrimPrefix(baseDomain, "www.")

	if baseDomain == "" {
		log.Printf("Failed to extract valid base domain from input: %s", input)
		// Log event against the original input if base domain extraction failed
		db.LogEvent(input, "analysis_failed", map[string]string{
			"error": "Invalid domain or URL format: " + input,
		})
		return
	}

	log.Printf("Base domain for analysis: %s", baseDomain)

	// 2. Log analysis started for the base domain
	if err := db.LogEvent(baseDomain, "analysis_started", nil); err != nil {
		log.Printf("Error logging analysis start for base domain %s: %v", baseDomain, err)
		// Don't return here, try analysis anyway but log the error
	}

	// 3. Determine URLs to check
	urlsToCheck := []string{}
	// Always check the base domain first
	urlsToCheck = append(urlsToCheck, "https://"+baseDomain+"/")
	// If the original input included 'www.', check that variant first instead.
	// Or if the original input didn't have 'www.' but was a domain, add 'www.' variant.
	originalHost := ""
	if potentialURL != nil {
		originalHost = potentialURL.Hostname()
	} else {
		originalHost = input // Use input if parsing failed initially
	}

	if strings.HasPrefix(originalHost, "www.") {
		// Input was www.example.com, check it first
		urlsToCheck = []string{"https://www." + baseDomain + "/"}
		// Add base domain as the second option
		urlsToCheck = append(urlsToCheck, "https://"+baseDomain+"/")
	} else {
		// Input was example.com (or something else), check base first, then www
		// Base is already added, add www version
		urlsToCheck = append(urlsToCheck, "https://www."+baseDomain+"/")
	}

	// Remove duplicates just in case logic above created any
	uniqueURLs := make(map[string]struct{})
	finalURLs := []string{}
	for _, u := range urlsToCheck {
		if _, exists := uniqueURLs[u]; !exists {
			uniqueURLs[u] = struct{}{}
			finalURLs = append(finalURLs, u)
		}
	}

	log.Printf("URLs to check for base domain %s: %v", baseDomain, finalURLs)

	// 4. Perform analysis sequentially
	var firstErrorPayload map[string]string
	analysisSucceeded := false
	successPayload := map[string]string{}

	for _, currentURL := range finalURLs {
		isShopify, payload := performSingleAnalysis(currentURL)
		if isShopify {
			log.Printf("Analysis succeeded for base domain %s via URL %s", baseDomain, currentURL)
			analysisSucceeded = true
			successPayload = payload
			// Add the URL that succeeded to the payload
			successPayload["checked_url"] = currentURL
			break // Stop checking on first success
		} else {
			log.Printf("Analysis failed for URL %s (related to base domain %s)", currentURL, baseDomain)
			if firstErrorPayload == nil {
				// Store the error from the first attempt (usually the more canonical domain)
				firstErrorPayload = payload
				firstErrorPayload["checked_url"] = currentURL // Record which URL failed
			}
		}
	}

	// 5. Log final result for the base domain
	if analysisSucceeded {
		db.LogEvent(baseDomain, "analysis_succeeded", successPayload)
	} else {
		log.Printf("Analysis failed for base domain %s after checking all variants.", baseDomain)
		// Use the payload from the first failed check, or a generic one if somehow that's nil
		if firstErrorPayload == nil {
			firstErrorPayload = map[string]string{
				"error": "Analysis failed for all checked URLs.",
			}
		}
		// Add the base domain to the final error payload for clarity
		firstErrorPayload["base_domain"] = baseDomain
		db.LogEvent(baseDomain, "analysis_failed", firstErrorPayload)
	}
}
