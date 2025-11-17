package routes

import (
	"github.com/mukailasam/igo"
)

type handlerInterface interface {
	SignUp(c *igo.Context)
	LoginUser(c *igo.Context)
	CreateProject(c *igo.Context)
	StartTimerEntry(c *igo.Context)
	StopTimerEntry(c *igo.Context)
	ListTimerEntries(c *igo.Context)
	ListProjects(c *igo.Context)
}
