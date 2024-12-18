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

	// start server
	e.Logger.Fatal(e.Start(utils.GetPort()))
}
