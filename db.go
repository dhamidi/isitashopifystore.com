package main

import (
	"database/sql"
	"encoding/json"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the SQL database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}

	log.Println("Database connection established with WAL mode")
	return &Database{db: db}, nil
}

// CreateEventsTable ensures the events table exists
func (d *Database) CreateEventsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL,
		event_type TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		payload JSON
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Events table created or already exists")
	return nil
}

// LogEvent records an event in the database
func (d *Database) LogEvent(domain, eventType string, payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO events (domain, event_type, payload)
	VALUES (?, ?, ?)`

	_, err = d.db.Exec(query, domain, eventType, string(jsonPayload))
	if err != nil {
		return err
	}

	log.Printf("Event logged: domain=%s, type=%s", domain, eventType)
	return nil
}

// GetLatestAnalysisResult retrieves the most recent analysis result for a domain
func (d *Database) GetLatestAnalysisResult(domain string) (string, string, error) {
	var eventType, payload string
	err := d.db.QueryRow(`
		SELECT event_type, payload 
		FROM events 
		WHERE domain = ? 
		ORDER BY id DESC 
		LIMIT 1`, domain).Scan(&eventType, &payload)
	
	return eventType, payload, err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
