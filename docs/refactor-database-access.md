# Database Access Refactoring Plan

## Overview
Move all database operations into db.go to create a clean separation of concerns and make the database layer more maintainable.

## Current State Analysis
Currently we have these database operations:
1. initDB() - Creates database connection
2. createEventsTable() - Creates events table
3. logEvent() - Logs events to database
4. Direct SQL queries in handler.go for:
   - Checking analysis status
   - Reading analysis results

## Step-by-Step Refactoring Plan

### 1. Create Database Type
In db.go, define:
```go
type Database struct {
    db *sql.DB
}
```

### 2. Move Existing Functions
Convert existing functions to methods:
```go
func NewDatabase(dbPath string) (*Database, error)  // renamed from initDB
func (d *Database) CreateEventsTable() error       // renamed from createEventsTable
func (d *Database) LogEvent(domain, eventType string, payload interface{}) error
```

### 3. Add New Methods
Move queries from handler.go into new methods:
```go
func (d *Database) GetLatestAnalysisResult(domain string) (string, string, error)
```

### 4. Update Code Structure
1. Remove global `db` variable from handler.go
2. Update main.go to use Database type
3. Update handlers to use Database methods
4. Remove direct SQL queries from handlers

### 5. Testing Steps
1. Test each database method
2. Verify all handlers work with new database layer
3. Check error handling
4. Verify logging still works

## Implementation Order
1. Create Database type
2. Convert existing functions to methods
3. Add new methods for handler queries
4. Update handlers one at a time
5. Test thoroughly

## Success Criteria
- All database operations moved to db.go
- No SQL queries in handler.go
- All tests passing
- No regression in functionality
- Clean error handling
- Proper logging maintained
