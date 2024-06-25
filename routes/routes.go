package routes

import (
	"net/http"

	"library/controllers/AuthController"
	"library/controllers/AuthorController"
	"library/controllers/BookController"
	"library/controllers/BorrowController"
	"library/controllers/CategoryController"
	"library/controllers/InventoryController"
	"library/controllers/RoleController"
	"library/controllers/UserController"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "haii")
	})

	// Media
	e.Static("avatar", "assets/images/avatar")
	e.Static("thumbnail", "assets/images/thumbnail")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/auth/login", AuthController.Login)
	e.POST("/auth/register", AuthController.Register)

	// Middleware JWT

	// Group for authenticated routes
	authenticated := e.Group("")
	authMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	})

	authenticated.Use(authMiddleware)

	// Profile
	authenticated.GET("/auth/profile", AuthController.Profile)
	authenticated.POST("/auth/profile", AuthController.UpdateProfile)

	// Roles
	authenticated.GET("/roles", RoleController.GetAll)
	authenticated.POST("/roles", RoleController.Create)
	authenticated.GET("/roles/:id", RoleController.Show)
	authenticated.PUT("/roles/:id", RoleController.Update)
	authenticated.DELETE("/roles/:id", RoleController.Delete)

	// Users
	authenticated.GET("/users", UserController.GetAll)
	authenticated.POST("/users", UserController.Create)
	authenticated.GET("/users/:id", UserController.Show)
	authenticated.PUT("/users/:id", UserController.Update)

	//Author
	authenticated.GET("/author", AuthorController.GetAll)
	authenticated.POST("/author", AuthorController.Create)
	authenticated.PUT("/author/:id", AuthorController.Update)
	authenticated.GET("/author/:id", AuthorController.Show)
	authenticated.DELETE("/author/:id", AuthorController.Delete)

	//Category
	authenticated.GET("/category", CategoryController.GetAll)
	authenticated.POST("/category", CategoryController.Create)
	authenticated.PUT("/category/:id", CategoryController.Update)
	authenticated.GET("/category/:id", CategoryController.Show)
	authenticated.DELETE("/category/:id", CategoryController.Delete)

	//Books
	authenticated.GET("/book", BookController.GetAll)
	authenticated.POST("/book", BookController.Create)
	authenticated.PUT("/book/:id", BookController.Update)
	authenticated.GET("/book/:id", BookController.Show)
	authenticated.DELETE("/book/:id", BookController.Delete)
	authenticated.GET("/book/:id/inventory", BookController.ShowInvetory)

	//Inventory
	authenticated.GET("/inventory", InventoryController.GetAll)
	authenticated.POST("/inventory", InventoryController.Create)
	authenticated.PUT("/inventory", InventoryController.Update)

	//Books
	authenticated.GET("/borrow", BorrowController.GetAll)
	authenticated.GET("/borrow/:id", BorrowController.Show)
	authenticated.POST("/borrow", BorrowController.Create)
	authenticated.PUT("/borrow/:id", BorrowController.Update)
	authenticated.DELETE("/borrow/:id", BorrowController.Delete)

	authenticated.GET("/borrowing", BorrowController.HistoryBorrowing)
	authenticated.POST("/borrowing", BorrowController.BorrowingBook)
	authenticated.PUT("/borrowing", BorrowController.ReturnBorrow)

	return e
}
