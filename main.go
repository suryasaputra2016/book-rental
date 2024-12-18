package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suryasaputra2016/book-rental/config"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/service"
)

func main() {
	db := config.ConnectDB()
	userRepo := repo.NewUserRepo(db)
	userservice := service.NewUserService(userRepo)

	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "internal server error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message.(string)
		}
		c.JSON(code, map[string]string{
			"error": msg,
		})
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/register", userservice.CreateUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	err := e.Start(":" + port)
	if err != nil {
		e.Logger.Fatal(err)
	}
}
