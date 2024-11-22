package menu

import "github.com/labstack/echo/v5"

func RegisterMenuRoutes(g *echo.Group, menuService *MenuService) {

	u := g.Group("/menu")

	u.GET("", getAllHandler(menuService))
	u.POST("", createItemHandler(menuService))
}
