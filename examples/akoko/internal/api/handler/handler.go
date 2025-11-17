package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mukailasam/akoko/pkg/auth"
	"github.com/mukailasam/akoko/util"
	"github.com/mukailasam/igo"
)

type Handler struct {
	apiService apiServiceInterface
	logger     loggerInterface
}

type Response map[string]string

type contextKey string

var uidKey = contextKey("uid")

func NewHandler(apiService apiServiceInterface, logger loggerInterface) *Handler {
	return &Handler{
		apiService: apiService,
		logger:     logger,
	}
}

func (h *Handler) SignUp(c *igo.Context) {
	var req struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.logger.Error("badrequet", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
		return
	}

	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	err := util.IsEmpty(email, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{"Message": err.Error()})
		return
	}
	err = util.ValidateEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{"Message": err.Error()})
		return
	}
	err = util.ValidatePassword(password)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{"Message": err.Error()})
		return
	}

	err = h.apiService.CreateAccount(c.Request.Context(), email, password)
	if err != nil {
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{"Message": "Account created successfully"})
}

func (h *Handler) LoginUser(c *igo.Context) {
	var req struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.logger.Error("badrequet", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
		return
	}

	user, err := h.apiService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, map[string]string{"error_message": "invalid credentials"})
			return

		}
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) CreateProject(c *igo.Context) {
	uidVal := c.Request.Context().Value(uidKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		c.Abort(403, "unauthorized: UID missing")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Client string `json:"client"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.logger.Error("badrequet", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
		return
	}

	err := h.apiService.CreateProject(c.Request.Context(), uid, req.Name, req.Client)
	if err != nil {
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{"Message": "Project created successfully"})
}

func (h *Handler) StartTimerEntry(c *igo.Context) {
	uidVal := c.Request.Context().Value(uidKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		c.Abort(403, "unauthorized: UID missing")
		return
	}
	var req struct {
		ProjectID   string `json:"project_id"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.logger.Error("badrequet", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
		return
	}

	entry, err := h.apiService.StartTimerEntry(c.Request.Context(), uid, req.ProjectID, req.Description)
	if err != nil {
		if err.Error() == "active" {
			h.logger.Error("conflict", "Error", fmt.Errorf("there is active timer for user: %v", uid))
			c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
			return
		}

		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": "internal server"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *Handler) StopTimerEntry(c *igo.Context) {
	uidVal := c.Request.Context().Value(uidKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		c.Abort(403, "unauthorized: UID missing")
		return
	}
	var req struct {
		ProjectID   string `json:"project_id"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.logger.Error("badrequet", "Error", err)
		c.JSON(http.StatusBadRequest, Response{"Message": "bad request"})
		return
	}

	res, err := h.apiService.StartTimerEntry(c.Request.Context(), uid, req.ProjectID, req.Description)
	if err != nil {
		h.logger.Error("internal server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ListTimerEntries(c *igo.Context) {
	uidVal := c.Request.Context().Value(uidKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		c.Abort(403, "unauthorized: UID missing")
		return
	}
	res, err := h.apiService.ListTimerEntries(c.Request.Context(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{"Message": "no timer entries found"})
			return
		}
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ListProjects(c *igo.Context) {
	uidVal := c.Request.Context().Value(uidKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		c.Abort(403, "unauthorized: UID missing")
		return
	}
	res, err := h.apiService.ListProjects(c.Request.Context(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{"Message": "no project found"})
			return
		}
		h.logger.Error("interna server error", "Error", err)
		c.JSON(http.StatusInternalServerError, Response{"Message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}
