package model

import (
	"time"
)

var ActiveTimer bool

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Project struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Client    string    `db:"client" json:"client"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type TimeEntry struct {
	ID              string    `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"user_id"`
	ProjectID       string    `db:"project_id" json:"project_id"`
	Description     string    `db:"description" json:"description"`
	StartTime       time.Time `db:"start_time" json:"start_time"`
	EndTime         time.Time `db:"end_time" json:"end_time"`
	DurationSeconds int64     `db:"duration_seconds" json:"duration_seconds"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
