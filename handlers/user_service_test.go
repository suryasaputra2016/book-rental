package handlers

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/suryasaputra2016/book-rental/repo"
// )

// var userRepoMock = &repo.UserRepoMock{Mock: mock.Mock{}}
// var userServiceMock = userService{ur: userRepoMock}

// func TestLoginEmailNotFound(t *testing.T) {
// 	userRepoMock.Mock.On("FindUserByID", 1).Return(nil)

// 	var c echo.Context
// 	// c.SetRequest(&http.Request{
// 	// 	Body: ,
// 	// })
// 	err := (*userService).Login(c)
// 	assert.NotNil(t, err) // expected : error =  not nil
// 	assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, "email cannot be found"), err.Error(), "Test Error product not found")
// }
