package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func DBConnect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
	    uid TEXT PRIMARY KEY,
	    email TEXT NOT NULL UNIQUE,
	    password_hash TEXT NOT NULL,
	    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS projects (
	    id TEXT PRIMARY KEY,
	    uid TEXT NOT NULL,
	    name TEXT NOT NULL,
	    client TEXT,
	    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
	    FOREIGN KEY (uid) REFERENCES users(uid) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS time_entries (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    uid TEXT NOT NULL,
	    pid TEXT NOT NULL,
	    description TEXT,
	    start_time DATETIME NOT NULL,
	    end_time DATETIME,
	    duration_seconds INTEGER,
	    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
	    FOREIGN KEY (uid) REFERENCES users(uid) ON DELETE CASCADE,
	    FOREIGN KEY (pid) REFERENCES projects(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_time_entries_user_start ON time_entries(uid, start_time);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
