package main

import (
	"database/sql"
	"encoding/json"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "isitashopifystore.db")
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
	return db, nil
}

func createEventsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL,
		event_type TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		payload JSON
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Events table created or already exists")
	return nil
}

func logEvent(db *sql.DB, domain, eventType string, payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO events (domain, event_type, payload)
	VALUES (?, ?, ?)`

	_, err = db.Exec(query, domain, eventType, string(jsonPayload))
	if err != nil {
		return err
	}

	log.Printf("Event logged: domain=%s, type=%s", domain, eventType)
	return nil
} 