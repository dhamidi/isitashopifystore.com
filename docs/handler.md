# handler.go

## Responsibility
The `handler.go` file defines HTTP handlers for the application's routes, managing web requests and responses for the Shopify store detection service.

## Structure

1. **Package Declaration**: Part of the `main` package.

2. **Import Section**: Imports packages for database access, JSON handling, HTML templates, logging, HTTP, URL parsing, and string manipulation.

3. **Global Variables**:
   - `db`: A global database connection variable.

4. **Type Definitions**:
   - `AnalysisResult`: Struct for storing and transmitting analysis results.

5. **HTTP Handlers**:
   - `landingPageHandler`: Handles GET requests to the landing page and processes URL form submissions. On POST, it extracts the domain, attempts to normalize it, and redirects to the result page path using the *original* domain input.
   - `resultPageHandler`: Handles requests like `/{domain}`. Extracts the domain from the path, normalizes it to the *base domain* (removes `www.`), and queries the database using the *base domain*. If no result is found, it triggers a background analysis using the *original* domain input and shows the polling page. If a result is found, it renders the result page (`YES`/`NO`) based on the stored `analysis_succeeded` or `analysis_failed` status.
   - `statusHandler`: Handles requests like `/status/{domain}`. Extracts the domain from the path, normalizes it to the *base domain*, and queries the database for the latest event using the *base domain*. Returns a JSON response indicating the status (`in_progress`, `succeeded`, `failed`). If no analysis is found for the base domain, it triggers a new background analysis using the *original* domain input.
   - `faviconHandler`: Serves the application's favicon.

6. **Helper Functions**:
   - `isValidDomain`: Validates domain name format.

## Opportunities for Abstraction

1. **Handler Organization**: Move handlers to separate files based on responsibility.

2. **Input Validation**: Create a dedicated validation layer instead of inline validation in handlers.

3. **Response Formatting**: Abstract the common patterns for JSON and HTML responses.

4. **Database Interaction**: Move the database queries to a dedicated data access layer.

5. **Error Handling**: Implement consistent error handling and reporting across all handlers.

6. **URL Processing**: Extract the common URL/domain parsing logic into a shared utility. 