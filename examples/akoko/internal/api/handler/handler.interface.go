package handler

import (
	"context"

	"github.com/mukailasam/akoko/internal/database/repository/model"
)

type apiServiceInterface interface {
	CreateAccount(ctx context.Context, email, password string) error
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	CreateProject(ctx context.Context, uid, name, client string) error
	StartTimerEntry(ctx context.Context, uid, pid, description string) (*model.TimeEntry, error)
	StopTimerEntry(ctx context.Context, uid, pid, name string) (*model.TimeEntry, error)
	ListProjects(ctx context.Context, uid string) ([]*model.Project, error)
	ListTimerEntries(ctx context.Context, uid string) ([]*model.TimeEntry, error)
}

type loggerInterface interface {
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
