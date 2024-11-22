package menu

import (
	"github.com/labstack/echo/v5"
	"net/http"
)

func getAllHandler(menuService *MenuService) echo.HandlerFunc {
	return func(c echo.Context) error {
		menuService.getAllMenuItems()
		return c.String(http.StatusOK, "User registered successfully")
	}
}

func createItemHandler(menuService *MenuService) echo.HandlerFunc {
	return func(c echo.Context) error {
		menuService.createItem()
		return c.String(http.StatusOK, "User registered successfully")
	}
}
