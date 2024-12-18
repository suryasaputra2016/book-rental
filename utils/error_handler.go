package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// http error handler function
func HTTPErrorHandler(err error, c echo.Context) {
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
