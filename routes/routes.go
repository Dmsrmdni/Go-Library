package routes

import (
	"net/http"

	"library/controllers/AuthorController"
	"library/controllers/BookController"
	"library/controllers/CategoryController"
	"library/controllers/RoleController"
	"library/controllers/UserController"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "haii")
	})

	// Media
	e.Static("avatar", "assets/images/avatar")
	e.Static("thumbnail", "assets/images/thumbnail")

	// Roles
	e.GET("/roles", RoleController.GetAll)
	e.POST("/roles", RoleController.Create)
	e.GET("/roles/:id", RoleController.Show)
	e.PUT("/roles/:id", RoleController.Update)
	e.DELETE("/roles/:id", RoleController.Delete)

	// Users
	e.GET("/users", UserController.GetAll)
	e.POST("/users", UserController.Create)
	e.GET("/users/:id", UserController.Show)
	e.PUT("/users/:id", UserController.Update)

	//Author
	e.GET("/author", AuthorController.GetAll)
	e.POST("/author", AuthorController.Create)
	e.PUT("/author/:id", AuthorController.Update)
	e.GET("/author/:id", AuthorController.Show)
	e.DELETE("/author/:id", AuthorController.Delete)

	//Category
	e.GET("/category", CategoryController.GetAll)
	e.POST("/category", CategoryController.Create)
	e.PUT("/category/:id", CategoryController.Update)
	e.GET("/category/:id", CategoryController.Show)
	e.DELETE("/category/:id", CategoryController.Delete)

	//Books
	e.GET("/book", BookController.GetAll)
	e.POST("/book", BookController.Create)
	e.PUT("/book/:id", BookController.Update)
	e.GET("/book/:id", BookController.Show)
	e.DELETE("/book/:id", BookController.Delete)

	return e
}
