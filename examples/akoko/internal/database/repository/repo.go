package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/mukailasam/akoko/internal/database/repository/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db DatabaseInterface
}

func NewRepository(db DatabaseInterface) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	uid := uuid.New().String()

	query := `INSERT INTO users (uid, email, password_hash) VALUES (?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, uid, email, string(hash))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := `SELECT uid, email, password_hash, created_at FROM users WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *Repository) CreateProject(ctx context.Context, uid, name, client string) error {
	query := `INSERT INTO projects (pid, uid, name, client) VALUES (?, ?, ?, ?)`
	pid := uuid.New().String()

	_, err := r.db.ExecContext(ctx, query, pid, uid, name, client)
	if err != nil {
		return err
	}

	return nil

}

func (r *Repository) ListProjects(ctx context.Context, uid string) ([]*model.Project, error) {
	query := `SELECT id, uid, name, client, created_at FROM projects WHERE uid = ?`

	rows, err := r.db.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*model.Project{}
	for rows.Next() {
		project := &model.Project{}
		rows.Scan(&project.ID, &project.UserID, &project.Name, &project.Client, &project.CreatedAt)
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *Repository) IsTimerActive(ctx context.Context, uid string) error {
	var existingID string
	var err error

	query := `SELECT id FROM time_entries WHERE uid = ? AND end_time IS NULL`

	row := r.db.QueryRowContext(ctx, query, uid)
	if err := row.Scan(existingID); err != nil {
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}

	model.ActiveTimer = true

	return nil
}

func (r *Repository) StartTimerEntry(ctx context.Context, uid, pID, description string) (*model.TimeEntry, error) {
	query := `INSERT INTO time_entries (id, uid, pid, description, start_time) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`

	timeEntryID := uuid.New().String()

	_, err := r.db.ExecContext(ctx, query, timeEntryID, uid, pID, description)
	if err != nil {
		return nil, err
	}

	res, err := r.GetTimerEntry(ctx, timeEntryID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) StopTimerEntry(ctx context.Context, uid, pid, timeEntryID string) (*model.TimeEntry, error) {
	query := `UPDATE time_entries SET end_time = CURRENT_TIMESTAMP WHERE uid = ? AND id = ? AND end_time IS NULL`

	_, err := r.db.ExecContext(ctx, query, uid, timeEntryID)
	if err != nil {
		return nil, err
	}

	res, err := r.GetTimerEntry(ctx, timeEntryID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) GetTimerEntry(ctx context.Context, timeEntryID string) (*model.TimeEntry, error) {
	query := `SELECT id, uid, pid, description, start_time, end_time, duration_seconds, created_at FROM time_entries WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, timeEntryID)

	te := &model.TimeEntry{}
	err := row.Scan(&te.ID, &te.UserID, &te.ProjectID, &te.Description, &te.StartTime, &te.EndTime, &te.DurationSeconds, &te.CreatedAt)
	if err != nil {
		return nil, err
	}
	return te, nil
}

func (r *Repository) ListTimerEntries(ctx context.Context, uid string) ([]*model.TimeEntry, error) {
	query := `SELECT id, uid, project_id, description, start_time, end_time, duration_seconds, created_at FROM time_entries WHERE uid = ? ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*model.TimeEntry
	for rows.Next() {
		te := &model.TimeEntry{}
		rows.Scan(&te.ID, &te.ID, &te.ProjectID, &te.Description, &te.StartTime, &te.EndTime, &te.DurationSeconds, &te.CreatedAt)
		entries = append(entries, te)
	}

	return entries, nil
}
