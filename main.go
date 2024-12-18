package main

import (
	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/config"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/service"
)

func main() {
	db := config.ConnectDB()
	userRepo := repo.NewUserRepo(db)
	userservice := service.NewUserService(userRepo)

	e := echo.New()
	e.GET("/register", userservice.CreateUser)
}
