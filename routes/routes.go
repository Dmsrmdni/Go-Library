package routes

import (
	"net/http"

	"library/controllers/RoleController"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "haii")
	})

	// Roles
	e.GET("/roles", RoleController.GetAll)
	e.POST("/roles", RoleController.Create)
	e.GET("/roles/:id", RoleController.Show)
	e.PUT("/roles/:id", RoleController.Update)
	e.DELETE("/roles/:id", RoleController.Delete)

	return e
}
