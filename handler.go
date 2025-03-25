package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var db *sql.DB

type AnalysisResult struct {
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	IsShopify bool `json:"is_shopify"`
	Domain string `json:"domain,omitempty"`
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
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Is it a Shopify Store?</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f5f5f5; }
        .container { text-align: center; padding: 2rem; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); max-width: 500px; width: 90%; }
        form { display: flex; gap: 0.5rem; margin-top: 1rem; }
        input[type="url"] { flex: 1; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; font-size: 1rem; }
        button { padding: 0.5rem 1rem; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 1rem; }
        button:hover { background: #2980b9; }
        .error { color: #e74c3c; margin-top: 0.5rem; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Is it a Shopify Store?</h1>
        <p>Enter a URL to check if it's a Shopify store</p>
        <form method="POST">
            <input type="url" name="url" placeholder="https://example.com" required>
            <button type="submit">Check</button>
        </form>
    </div>
</body>
</html>`
	t := template.Must(template.New("landing").Parse(tmpl))
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
	err := db.QueryRow(`
		SELECT event_type, payload 
		FROM events 
		WHERE domain = ? 
		ORDER BY id DESC 
		LIMIT 1`, domain).Scan(&result.Status, &result.Reason)
	
	if err == sql.ErrNoRows {
		log.Printf("No analysis found for domain: %s, starting new analysis", domain)
		// No analysis exists, trigger background analysis
		go analyzeDomain(domain)
		
		// Show polling page
		log.Printf("Rendering polling page for domain: %s", domain)
		tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Analyzing {{.Domain}}</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f5f5f5; }
        .container { text-align: center; padding: 2rem; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .spinner { width: 40px; height: 40px; border: 4px solid #f3f3f3; border-top: 4px solid #3498db; border-radius: 50%; animation: spin 1s linear infinite; margin: 20px auto; }
        @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
    </style>
</head>
<body>
    <div class="container">
        <h1>Analyzing {{.Domain}}</h1>
        <div class="spinner"></div>
        <p>Checking if this is a Shopify store...</p>
    </div>
    <script>
        function checkStatus() {
            fetch('/status/{{.Domain}}')
                .then(response => response.json())
                .then(data => {
                    if (data.status === 'succeeded' || data.status === 'failed') {
                        window.location.reload();
                    } else {
                        setTimeout(checkStatus, 1000);
                    }
                });
        }
        checkStatus();
    </script>
</body>
</html>`
		t := template.Must(template.New("polling").Parse(tmpl))
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
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Is {{.Domain}} a Shopify Store?</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f5f5f5; }
        .container { text-align: center; padding: 2rem; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .result { font-size: 4rem; font-weight: bold; margin: 1rem 0; }
        .yes { color: #2ecc71; }
        .no { color: #e74c3c; }
        .reason { color: #666; margin-top: 1rem; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Is {{.Domain}} a Shopify Store?</h1>
        <div class="result {{if .IsShopify}}yes{{else}}no{{end}}">
            {{if .IsShopify}}YES{{else}}NO{{end}}
        </div>
        {{if .IsShopify}}
        <div class="reason">
            {{.Reason}}
        </div>
        {{end}}
    </div>
</body>
</html>`
	t := template.Must(template.New("result").Parse(tmpl))
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
	err := db.QueryRow(`
		SELECT event_type, payload 
		FROM events 
		WHERE domain = ? 
		ORDER BY id DESC 
		LIMIT 1`, domain).Scan(&result.Status, &result.Reason)
	
	if err == sql.ErrNoRows {
		log.Printf("No analysis found for domain in status check: %s", domain)
		// No analysis exists yet
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

// isValidDomain checks if a string is a valid domain name
func isValidDomain(domain string) bool {
	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$`)
	return domainRegex.MatchString(domain)
}