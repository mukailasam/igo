package api

import (
	"context"

	"github.com/mukailasam/akoko/internal/database/repository/model"
)

type repositoryInterface interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateProject(ctx context.Context, uid, name, client string) error
	ListProjects(ctx context.Context, uid string) ([]*model.Project, error)
	IsTimerActive(ctx context.Context, uid string) error
	StartTimerEntry(ctx context.Context, uid, pid, description string) (*model.TimeEntry, error)
	StopTimerEntry(ctx context.Context, uid, pid, timeEntryID string) (*model.TimeEntry, error)
	GetTimerEntry(ctx context.Context, timeEntryID string) (*model.TimeEntry, error)
	ListTimerEntries(ctx context.Context, uid string) ([]*model.TimeEntry, error)
}

type loggerInterface interface {
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
