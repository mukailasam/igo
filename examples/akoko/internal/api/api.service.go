package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/mukailasam/akoko/internal/database/repository/model"
	"golang.org/x/crypto/bcrypt"
)

type ApiService struct {
	repo   repositoryInterface
	logger loggerInterface
}

func NewAPIService(repo repositoryInterface, logger loggerInterface) *ApiService {
	return &ApiService{
		repo:   repo,
		logger: logger,
	}
}

func (api *ApiService) CreateAccount(ctx context.Context, email, password string) error {
	err := api.repo.CreateUser(ctx, email, password)
	if err != nil {
		api.logger.Error(err.Error())
		return err
	}

	return nil
}

func (api *ApiService) LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	user, err := api.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (api *ApiService) CreateProject(ctx context.Context, uid, name, client string) error {
	err := api.repo.CreateProject(ctx, uid, name, client)
	if err != nil {
		api.logger.Error(err.Error())
		return err
	}

	return nil
}

func (api *ApiService) StartTimerEntry(ctx context.Context, uid, pid, description string) (*model.TimeEntry, error) {
	err := api.repo.IsTimerActive(ctx, uid)
	if err != nil && !model.ActiveTimer {
		return nil, fmt.Errorf("active")
	}

	te, err := api.repo.StartTimerEntry(ctx, uid, pid, description)
	if err != nil {
		api.logger.Error(err.Error())
		return nil, err
	}

	return te, nil
}

func (api *ApiService) StopTimerEntry(ctx context.Context, uid, pid, timeEntryID string) (*model.TimeEntry, error) {
	res, err := api.repo.StopTimerEntry(ctx, uid, pid, timeEntryID)
	if err != nil {
		api.logger.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func (api *ApiService) ListProjects(ctx context.Context, uid string) ([]*model.Project, error) {
	res, err := api.repo.ListProjects(ctx, uid)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *ApiService) ListTimerEntries(ctx context.Context, uid string) ([]*model.TimeEntry, error) {
	res, err := api.repo.ListTimerEntries(ctx, uid)
	if err != nil {
		return nil, err
	}

	return res, nil
}
