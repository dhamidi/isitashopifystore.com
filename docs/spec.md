# Project Overview

Create a website (isitashopifystore.com) with a single URL input form. When a user submits a URL, the backend checks if the domain has been analyzed before. If yes, it immediately shows a full-screen "yes" or "no" based on the analysis. If not, it displays an “analysis in progress” message, starts the analysis, and the frontend polls every second until the result is ready.

# Tech Stack

- **Language:** Go
- **Database:** SQLite
- **Frontend:** Rendered server-side by the Go backend

# Architecture and Flow

1. **User Input & Redirection:**

   - The landing page hosts a simple form with a URL input.
   - Upon submission, the backend extracts the domain and redirects the user to `isitashopifystore.com/domain.of.the.store.com`.

2. **Result Page Behavior:**

   - The page initially checks if an analysis exists for the domain.
   - If found, it renders a full-screen page with a "yes" or "no" based on the analysis.
   - If not, it shows an “analysis in progress” message and triggers frontend polling (once a second) to fetch updates.

3. **Backend Analysis Logic:**

   - **Pre-check:** Verify if the domain is already in the database.
   - **Start Analysis:**
     - Record an "analysis started" event in SQLite with a timestamp and any event details in the JSON payload.
     - Make an HTTP GET request to the provided URL, following redirects.
     - Upon reaching the final URL, check if the response code is 200.
   - **Analysis Decision:**
     - **Success:** If the response body includes “myshopify” or “cdn.shopify.com”, record an "analysis succeeded" event, including the matching string as the reason. Render a full-screen "yes" result.
     - **Failure:** If the page cannot be reached or does not include the required keywords, record an "analysis failed" event with details in the JSON payload. Render a full-screen "no" result.

4. **Frontend Polling:**
   - The frontend polls the server every second via AJAX or a similar mechanism.
   - Polling continues until the server indicates that analysis is complete (success or failure).

# Database Schema

Create a SQLite table named `events` with the following columns:

- **id** (INTEGER, primary key, auto-increment)
- **domain** (TEXT) – The domain being analyzed
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
