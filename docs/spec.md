# Project Overview

Create a website (isitashopifystore.com) with a single URL input form. When a user submits a URL, the backend checks if the domain has been analyzed before. If yes, it immediately shows a full-screen "yes" or "no" based on the analysis. If not, it displays an "analysis in progress" message, starts the analysis, and the frontend polls every second until the result is ready.

# Tech Stack

- **Language:** Go
- **Database:** SQLite
- **Frontend:** Rendered server-side by the Go backend

# Architecture and Flow

1. **User Input & Redirection:**

   - The landing page hosts a simple form with a URL input.
   - Upon submission, the backend extracts the *base domain* (e.g., `example.com` from `www.example.com`) and redirects the user to the path corresponding to the *original* input domain (e.g., `isitashopifystore.com/www.example.com`).

2. **Result Page (`/{domain}`):**

   - The handler extracts the domain from the request path (e.g., `www.example.com`).
   - It normalizes this to the *base domain* (e.g., `example.com`).
   - It checks the database for the latest analysis result associated with the *base domain*.
   - **If found:** Renders a full-screen page with "yes" or "no" based on the stored result.
   - **If not found:** 
     - Records an "analysis started" event for the *base domain*.
     - Triggers a background analysis process using the *original* input domain (`www.example.com`).
     - Renders an "analysis in progress" page.

3. **Backend Analysis Logic (`analyzeDomain` function):**

   - Receives the *original* input domain (e.g., `www.example.com`).
   - Extracts the *base domain* (e.g., `example.com`).
   - Logs an `analysis_started` event for the *base domain*.
   - Determines the URLs to check: It will attempt analysis on both the base domain (`https://example.com/`) and the `www.` prefixed domain (`https://www.example.com/`). The order depends on the original input.
   - **For each URL:** It calls a helper function (`performSingleAnalysis`) that makes an HTTP GET request, follows redirects (up to 10), and checks the response body for Shopify indicators ("myshopify", "cdn.shopify.com") or checks the checkout page (`https://checkout.{domain}/...`) for Shopify headers (`x-shopid`, `Server: Shopify`).
   - **Analysis Decision:**
     - **Success:** If *any* of the checked URLs show Shopify indicators, it records a single "analysis succeeded" event for the *base domain*, including details from the first successful check (like the indicator found and the URL checked) in the payload.
     - **Failure:** If *neither* URL check succeeds, it records a single "analysis failed" event for the *base domain*, including details from the first failed check in the payload.

4. **Frontend Polling (`/status/{domain}`):**
   - The frontend polls the status endpoint using the *original* domain (e.g., `/status/www.example.com`).
   - The status handler extracts the domain from the path and normalizes it to the *base domain* (`example.com`).
   - It queries the database for the latest event associated with the *base domain*.
   - It returns a JSON response indicating the status (`in_progress`, `succeeded`, `failed`).
   - Polling continues (every second) until the status is `succeeded` or `failed`, at which point the frontend reloads the result page (`/{domain}`).

# Database Schema

Create a SQLite table named `events` with the following columns:

- **id** (INTEGER, primary key, auto-increment)
- **domain** (TEXT) – The *base* domain being analyzed (e.g., `example.com`)
- **event_type** (TEXT) – One of: "analysis started", "analysis succeeded", "analysis failed"
- **timestamp** (DATETIME) – The time when the event occurred
- **payload** (JSON) – Event-specific details (e.g., the matching string found, error messages)

# Additional Details

- **Concurrency & Duplication:**

  - If multiple requests for the same domain arrive, ensure only one analysis runs at a time.
  - Subsequent requests should either wait for the current analysis or use the cached result.

- **Error Handling:**

  - No retry logic is needed.
  - In case of errors (e.g., network issues, non-200 responses), record the "analysis failed" event with error details in the payload.

- **Logging & Monitoring:**

  - Besides recording events in SQLite, consider logging key steps (optional) for debugging purposes.

- **Deployment Considerations:**
  - Ensure proper handling of SQLite file permissions and concurrency, especially under load.
  - Security: Validate and sanitize user input to avoid injection attacks.
