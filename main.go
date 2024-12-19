package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suryasaputra2016/book-rental/config"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/service"
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

// @license.name The MIT License
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://*****.herokuapp.com
// @BasePath /
func main() {
	// configure database, user repo, and user service
	db := config.ConnectDB()
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	rentalRepo := repo.NewRentRepo(db)
	rentService := service.NewRentService(rentalRepo)
	bookRepo := repo.NewBookRepo(db)

	bookService := service.NewBookService(bookRepo, userRepo, rentalRepo)

	e := echo.New()

	// use swaggo
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// error handler
	e.HTTPErrorHandler = utils.HTTPErrorHandler

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.POST("/register", userService.CreateUser)
	e.POST("/login", userService.Login)
	e.PUT("/topup", userService.Topup, middlewares.Authorization())
	e.GET("/rents", rentService.ShowRents, middlewares.Authorization())
	e.GET("/books", bookService.ShowBooks)
	e.POST("/books/rent", bookService.RentBook, middlewares.Authorization())
	e.POST("/books/return", bookService.ReturnBook, middlewares.Authorization())

	// start server
	e.Logger.Fatal(e.Start(utils.GetPort()))
}
