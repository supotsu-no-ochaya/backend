package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func TestHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		records, err := app.FindAllRecords("users", nil)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to fetch users"})
		}
		return e.JSON(http.StatusOK, records)
	}
}
