package user

import "github.com/labstack/echo/v5"

func RegisterUserRoutes(g *echo.Group, userService *UserService) {

	u := g.Group("/user")

	u.GET("", userRegisterHandler(userService))
}
