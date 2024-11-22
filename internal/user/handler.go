package user

import (
	"github.com/labstack/echo/v5"
	"net/http"
)

func userRegisterHandler(userService *UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		userService.register(username)
		return c.String(http.StatusOK, "User registered successfully")
	}
}
