package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var db *sql.DB

type AnalysisResult struct {
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	IsShopify bool `json:"is_shopify"`
}

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

func resultPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract domain from path
	domain := strings.TrimPrefix(r.URL.Path, "/")
	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	// Check if analysis exists
	var result AnalysisResult
	err := db.QueryRow(`
		SELECT event_type, payload 
		FROM events 
		WHERE domain = ? 
		ORDER BY timestamp DESC 
		LIMIT 1`, domain).Scan(&result.Status, &result.Reason)
	
	if err == sql.ErrNoRows {
		// No analysis exists, show polling page
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
	domain := strings.TrimPrefix(r.URL.Path, "/status/")
	if domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	// Set JSON content type
	w.Header().Set("Content-Type", "application/json")

	// Check if analysis exists
	var result AnalysisResult
	err := db.QueryRow(`
		SELECT event_type, payload 
		FROM events 
		WHERE domain = ? 
		ORDER BY timestamp DESC 
		LIMIT 1`, domain).Scan(&result.Status, &result.Reason)
	
	if err == sql.ErrNoRows {
		// No analysis exists yet
		json.NewEncoder(w).Encode(AnalysisResult{
			Status: "in_progress",
		})
		return
	}

	if err != nil {
		log.Printf("Error checking analysis status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Map event types to status
	switch result.Status {
	case "analysis_started":
		result.Status = "in_progress"
	case "analysis_succeeded":
		result.Status = "succeeded"
		result.IsShopify = true
	case "analysis_failed":
		result.Status = "failed"
	default:
		result.Status = "in_progress"
	}

	// Return the result
	json.NewEncoder(w).Encode(result)
} 