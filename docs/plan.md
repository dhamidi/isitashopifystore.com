# Step-by-Step Implementation Plan for isitashopifystore.com

This document outlines a granular plan to build a Go-based website that checks if a given URL is a Shopify store. The app will consist of a landing page with a URL form, a results page that shows "yes" or "no" based on the analysis, background analysis processing, and a SQLite-backed event log. All code will reside in a single flat Go package.

---

## Task 1: Project Setup

### 1.1 Create the Project Structure - DONE

- **Task:** Create a new directory (e.g., `isitashopifystore`) and add these files:
  - `main.go`
  - `db.go`
  - `handler.go`
  - `analysis.go`
- **Definition of Done:** The project directory exists with the listed files; each file starts with `package main` and is empty or contains a basic skeleton.
- **LLM Prompt:** Create a basic Go project structure in a new directory named 'isitashopifystore' with the files main.go, db.go, handler.go, and analysis.go. Each file should have 'package main' as the first line.

### 1.2 Initialize a Go Module - DONE

- **Task:** Run `go mod init isitashopifystore` to create the `go.mod` file.
- **Definition of Done:** A `go.mod` file exists with the module name `isitashopifystore`.
- **LLM Prompt:** Help me initialize a new Go module for the project named 'isitashopifystore' and generate the go.mod file.

## Task 2: Database Setup (db.go)

### 2.1 Setup SQLite Connection - DONE

- **Task:** Write a function to open the SQLite database file (`isitashopifystore.db`).
- **Definition of Done:** A function (e.g., `initDB()`) returns a valid `*sql.DB` connection.
- **LLM Prompt:** In db.go, write a Go function 'initDB()' that opens a SQLite database file named 'isitashopifystore.db' using the database/sql package and the go-sqlite3 driver, then returns the connection.

### 2.2 Create the Events Table - DONE

- **Task:** In the same file, write code to create an `events` table (if it doesn't exist) with these columns:
- `id` (INTEGER, primary key, auto-increment)
- `domain` (TEXT)
- `event_type` (TEXT)
- `timestamp` (DATETIME)
- `payload` (JSON)
- **Definition of Done:** The database creates the `events` table on startup if it doesn't already exist.
- **LLM Prompt:** In db.go, write a function that creates an 'events' table with columns 'id' (INTEGER, primary key, auto-increment), 'domain' (TEXT), 'event_type' (TEXT), 'timestamp' (DATETIME), and 'payload' (JSON). Ensure the table is created if it does not already exist.

### 2.3 Add an Event Logging Helper - DONE

- **Task:** Write a helper function `logEvent(domain, eventType string, payload interface{})` to insert events into the table.
- **Definition of Done:** Calling `logEvent` successfully inserts an event into the database.
- **LLM Prompt:** In db.go, add a function 'logEvent(domain, eventType string, payload interface{})' that inserts a new event into the 'events' table with the current timestamp. Use proper error handling.

## Task 3: HTTP Server Setup (main.go)

### 3.1 Configure the HTTP Server - DONE

- **Task:** Set up an HTTP server using `net/http` that routes requests to handlers.
- **Steps:**
- Import necessary packages.
- Call the DB initialization function from Task 2.
- Define routes for `/` (landing page) and `/status/` (polling endpoint) and `/[domain]` (results page).
- **Definition of Done:** The HTTP server starts on a designated port with the correct routes set up.
- **LLM Prompt:** In main.go, write a main() function that initializes the database (calling initDB()), sets up HTTP routes for '/' (landing page), '/status/' (polling), and '/[domain]' (results page), and starts the server using net/http.

## Task 4: Landing Page Implementation (handler.go)

### 4.1 Create the Landing Page Handler - DONE

- **Task:** Implement `landingPageHandler` to render an HTML form with a URL input.
- **Steps:**
- Write HTML that includes a form with a URL input field and a submit button.
- Use the POST method.
- **Definition of Done:** Visiting `/` displays a page with the URL form.
- **LLM Prompt:** In handler.go, create a function 'landingPageHandler' that writes an HTML page containing a form with a URL input field (using POST). The form should be simple and self-contained.

### 4.2 Process Form Submission and Redirect - DONE

- **Task:** Update the form handler to:
- Parse the submitted URL.
- Extract the domain.
- Redirect the user to `/{domain}`.
- **Definition of Done:** Submitting the form redirects to a URL based on the extracted domain.
- **LLM Prompt:** In handler.go, modify or add a function to handle form submissions from the landing page. It should extract the domain from the provided URL and then redirect the client to '/{domain}'.

## Task 5: Result Page & Polling Implementation (handler.go)

### 5.1 Create the Result Page Handler - DONE

- **Task:** Implement `resultPageHandler` that:
- Checks if an analysis result for the domain exists in the database.
- If complete, renders a full-screen "yes" or "no" (including the reason on success).
- If not complete, renders an "analysis in progress" message with embedded JavaScript that polls every second.
- **Definition of Done:** Visiting `/{domain}` shows either the final result or a polling page while analysis is in progress.
- **LLM Prompt:** In handler.go, implement a 'resultPageHandler' function that, for a given domain, checks if an analysis result exists. If it does, render a full-screen page showing 'yes' or 'no' (with the reason on success). If not, display an 'analysis in progress' message with JavaScript that polls the server every second.

### 5.2 Create the Polling Endpoint - DONE

- **Task:** Implement `statusHandler` to:
- Return the current analysis status as JSON.
- Include status values such as "in progress", "succeeded", "failed", and a reason if available.
- **Definition of Done:** The `/status/{domain}` endpoint returns accurate JSON data reflecting the analysis state.
- **LLM Prompt:** In handler.go, implement a 'statusHandler' that responds to '/status/{domain}' requests. It should return JSON indicating whether the analysis is 'in progress', 'succeeded', or 'failed', and include the reason if applicable.

## Task 6: Analysis Logic Implementation (analysis.go)

### 6.1 Implement the Domain Analysis Function - DONE

- **Task:** In `analysis.go`, write `analyzeDomain(domain string)` that:
- Logs an "analysis started" event.
- Makes an HTTP GET request to the given domain, following redirects.
- Verifies that the final response is 200.
- Searches the response body for "myshopify" or "cdn.shopify.com".
- Logs an "analysis succeeded" event with the found keyword as the reason if the check passes.
- Logs an "analysis failed" event with error details otherwise.
- **Definition of Done:** The function completes the analysis and logs the appropriate event in the database.
- **LLM Prompt:** In analysis.go, implement a function 'analyzeDomain(domain string)' that:

Logs an 'analysis started' event.

Performs an HTTP GET request to the given domain (following redirects).

Checks if the final response is 200.

Searches the response body for 'myshopify' or 'cdn.shopify.com'.

Logs an 'analysis succeeded' event with the matching keyword as the reason if found; otherwise, logs an 'analysis failed' event.

### 6.2 Integrate Background Analysis - DONE

- **Task:** Modify `resultPageHandler` to:
- Check if the analysis for the domain exists.
- If not, trigger `analyzeDomain(domain)` in a background goroutine.
- **Definition of Done:** When a user visits `/{domain}` for the first time, the analysis function is called asynchronously, and subsequent polling returns the updated status.
- **LLM Prompt:** In handler.go, update the 'resultPageHandler' so that if no analysis exists for the domain, it triggers 'analyzeDomain(domain)' in a separate goroutine to process the analysis in the background.

## Task 7: Input Validation & Logging Enhancements

### 7.1 Validate and Sanitize URL Input – DONE

- **Task:** Add code in the form submission handler to:
- Validate that the submitted URL is well-formed.
- If only a domain name is submitted, treat it as a https domain
- Sanitize the input to prevent injection attacks.
- **Definition of Done:** Only valid URLs are processed; invalid inputs result in a friendly error.
- **LLM Prompt:** In handler.go, add input validation and sanitization logic for the submitted URL in the form handler. Ensure the URL is well-formed and safe before processing.

### 7.2 Add Logging Statements – DONE

- **Task:** Throughout the code (main.go, handler.go, analysis.go), add logging for:
- Form submissions and redirections.
- Start and completion of the analysis.
- Any errors encountered.
- **Definition of Done:** Logs are printed to the console showing the flow of events and errors.
- **LLM Prompt:** In the project files, insert logging statements using Go's log package to track key actions such as form submissions, redirections, analysis start, analysis success, and analysis failures.

