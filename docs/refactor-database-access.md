# Database Access Refactoring Plan

## Overview
Move all database operations into db.go to create a clean separation of concerns and make the database layer more maintainable.

## Current State
- Database initialization in db.go
- Database operations scattered across handler.go
- Global db variable in handler.go

## Step-by-Step Refactoring Plan

### 1. Create Database Types
In db.go, define:
- `type Database struct { db *sql.DB }`
- `type AnalysisRecord struct` to match database schema

### 2. Create Database Methods
Add these functions to db.go:
```go
func NewDatabase(dbPath string) (*Database, error)
func (d *Database) Close() error
func (d *Database) GetAnalysisResult(domain string) (*AnalysisResult, error)
func (d *Database) SaveAnalysisResult(domain string, result *AnalysisResult) error
func (d *Database) UpdateAnalysisStatus(domain, status string) error
```

### 3. Update Database Schema
Add new table for analysis results:
```sql
CREATE TABLE IF NOT EXISTS analysis_results (
    domain TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    reason TEXT,
    is_shopify BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

### 4. Modify Existing Code
1. Remove global `db` variable from handler.go
2. Update main.go to use new Database type
3. Update all handlers to use Database methods
4. Remove direct SQL queries from handlers

### 5. Testing Steps
1. Create test database file
2. Test each new database method
3. Verify all handlers work with new database layer
4. Check error handling
5. Verify logging still works

### 6. Deployment Considerations
1. Create database migration script
2. Plan deployment sequence to avoid downtime
3. Update documentation

## Implementation Order
1. Create new types and interface
2. Implement database methods one at a time
3. Update handlers one at a time
4. Test thoroughly
5. Deploy with careful migration

## Success Criteria
- All database operations moved to db.go
- No SQL queries in handler.go
- All tests passing
- No regression in functionality
- Clean error handling
- Proper logging maintained
