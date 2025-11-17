package routes

import (
	"github.com/mukailasam/igo"
)

type Router struct {
	app     *igo.Igo
	handler handlerInterface
}

func NewRouter(app *igo.Igo, handler handlerInterface) *Router {
	return &Router{
		app:     app,
		handler: handler,
	}
}

func (r *Router) RegisterRoutes() {
	r.app.POST("/api/user", r.handler.SignUp)
	r.app.POST("/api/login", r.handler.LoginUser)
	r.app.POST("/api/project", r.handler.CreateProject)
	r.app.GET("/api/projects", r.handler.ListProjects)
	r.app.POST("/api/timer/start", r.handler.StartTimerEntry)
	r.app.POST("/api/timer/stop", r.handler.StopTimerEntry)
	r.app.GET("/api/timer", nil)
	r.app.GET("/api/timers", r.handler.ListTimerEntries)
}
