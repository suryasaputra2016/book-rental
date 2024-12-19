package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suryasaputra2016/book-rental/config"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/service"
	"github.com/suryasaputra2016/book-rental/utils"
)

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
	e.POST("/rentbook", bookService.RentBook, middlewares.Authorization())
	e.POST("/returnbook", bookService.ReturnBook, middlewares.Authorization())

	// start server
	e.Logger.Fatal(e.Start(utils.GetPort()))
}
