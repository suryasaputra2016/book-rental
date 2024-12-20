package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suryasaputra2016/book-rental/config"
	"github.com/suryasaputra2016/book-rental/handlers"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/services"
	"github.com/suryasaputra2016/book-rental/utils"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/suryasaputra2016/book-rental/docs"
)

// @title BookRental
// @version 1.0
// @description Hacktiv8 Phase 3 Mini Project
// @termsOfService http://swagger.io/terms/

// @contact.name Surya Saputra
// @contact.email sayasuryasaputra@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://*****.herokuapp.com
// @BasePath /
func main() {
	// database
	db := config.ConnectDB()

	// repositories
	userRepo := repo.NewUserRepo(db)
	rentalRepo := repo.NewRentRepo(db)
	bookRepo := repo.NewBookRepo(db)

	// handlers
	userHandler := handlers.NewUserHandler(services.NewUserService(userRepo))
	rentHandler := handlers.NewRentHandler(services.NewRentService(rentalRepo))
	bookHandler := handlers.NewBookHandler(services.NewBookService(userRepo, rentalRepo, bookRepo))

	e := echo.New()

	// use swaggo
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// error handler
	e.HTTPErrorHandler = utils.HTTPErrorHandler

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)
	e.PUT("/topup", userHandler.Topup, middlewares.Authorization())
	e.GET("/rents", rentHandler.ShowRents, middlewares.Authorization())
	e.GET("/books", bookHandler.ShowBooks)
	e.POST("/books/rent", bookHandler.RentABook, middlewares.Authorization())
	e.POST("/books/return", bookHandler.ReturnABook, middlewares.Authorization())

	// start server
	e.Logger.Fatal(e.Start(utils.GetPort()))
}
