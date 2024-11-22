package main

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/supotsu-no-ochaya/backend/internal/menu"
	"github.com/supotsu-no-ochaya/backend/internal/user"
	"log"
	"net/http"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		g := e.Router.Group("/cafe/api")

		userDB := user.UserDb{}
		userService := user.NewUserService(&userDB)
		user.RegisterUserRoutes(g, &userService)

		menuDB := menu.MenuDb{}
		menuService := menu.NewMenuService(&menuDB)
		menu.RegisterMenuRoutes(g, &menuService)

		e.Router.GET("/api/test", func(c echo.Context) error {
			// Fetch all users
			records, err := app.Dao().FindRecordsByExpr("users", nil)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to fetch users"})
			}

			return c.JSON(http.StatusOK, records)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
