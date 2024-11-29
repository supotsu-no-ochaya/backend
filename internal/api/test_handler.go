package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func TestHandler(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		records, err := app.Dao().FindRecordsByExpr("users", nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to fetch users"})
		}
		return c.JSON(http.StatusOK, records)
	}
}
