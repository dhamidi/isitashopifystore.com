# IsItAShopifyStore.com Project Documentation

## Project Overview
IsItAShopifyStore.com is a web application that analyzes websites to determine if they are built on the Shopify platform. The service accepts a URL or domain name, performs an analysis of the website's content, and reports whether it detects Shopify indicators.

## Project Structure

### Core Components

| File | Description | Documentation |
|------|-------------|---------------|
| [main.go](main.md) | Application entry point, server initialization, and route configuration | [Read more](main.md) |
| [handler.go](handler.md) | HTTP handlers for landing page, result page, and status endpoints | [Read more](handler.md) |
| [analysis.go](analysis.md) | Core logic for Shopify detection | [Read more](analysis.md) |
| [db.go](db.md) | Database initialization and event logging | [Read more](db.md) |

### Dependencies

The project uses the following main dependencies:
- Go standard library (net/http, database/sql, etc.)
- github.com/mattn/go-sqlite3 - SQLite database driver

### Application Flow

1. User submits a domain via the landing page form
2. System redirects to the domain-specific result page
3. If analysis doesn't exist, a background analysis is triggered
4. Client polls the status endpoint until analysis completes
5. System displays the analysis result when completed

### Database Schema

The application uses a SQLite database with a single `events` table:
- `id` - Unique identifier (auto-incremented)
- `domain` - The domain being analyzed
- `event_type` - Type of event (analysis_started, analysis_succeeded, analysis_failed)
- `timestamp` - When the event occurred
- `payload` - JSON data with additional information about the event

## Improvement Opportunities

Common themes for potential improvements across the codebase:

1. **Code Organization**: Separate responsibilities into dedicated packages
2. **Configuration Management**: Implement a configuration system
3. **Error Handling**: Develop a comprehensive error handling strategy
4. **Testing**: Add unit and integration tests
5. **Logging**: Improve logging with structured logs
6. **Performance**: Implement caching for analysis results
7. **Security**: Add rate limiting and input sanitization 