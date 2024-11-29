package routes

import (
	"github.com/pocketbase/pocketbase/core"

	"github.com/supotsu-no-ochaya/backend/internal/api"
)

func RegisterAPIRoutes(e *core.ServeEvent, app core.App) {
	apiGroup := e.Router.Group("/api")

	// /api/test endpoint
	apiGroup.GET("/test", api.TestHandler(app))

	// /api/export-json endpoint
	apiGroup.GET("/export-json", api.ExportJSONHandler(app))
}
