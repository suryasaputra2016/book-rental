package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middleware"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/utils"
)

type UserService interface {
	CreateUser(*entity.CreateUserRequest) error
	Login(*entity.LoginRequest) (string, error)
}

// user repository implementation with database connection
type userService struct {
	ur repo.UserRepo
}

func NewUserService(ur repo.UserRepo) *userService {
	return &userService{ur: ur}
}

func (us *userService) CreateUser(c echo.Context) error {
	// bind request body
	var req entity.CreateUserRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// email and password  validation
	if err := utils.IsEmailandPasswordFine(req.Email, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// check email if it already exists
	if _, err := us.ur.FindUserByEmail(req.Email); err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "email is already in use")
	}

	// hash password
	passwordHash, err := utils.GenerateHash(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't hash password")
	}

	//define new user
	newUser := entity.User{
		Email:         req.Email,
		PasswordHash:  passwordHash,
		DepositAmount: 0.0,
	}

	// add new user to database and return
	if us.ur.AddUser(&newUser) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	//define response
	res := &entity.CreateUserRepsonse{
		Message: "User successfully created",
		NewUserData: entity.CreateUserResponseData{
			Email:         newUser.Email,
			DepositAmount: newUser.DepositAmount,
		},
	}

	// return response
	return c.JSON(http.StatusCreated, *res)
}

func (us *userService) Login(reqPtr *entity.LoginRequest) (*string, error) {

	// check email and get member info
	var userPtr *entity.User
	userPtr, err := us.ur.FindUserByEmail(reqPtr.Email)
	if err != nil {
		return nil, errors.New("username cannot be found")
	}

	// check password
	err = utils.CheckPassword(reqPtr.Password, userPtr.PasswordHash)
	if err != nil {
		return nil, errors.New("password doesn't match")
	}

	// generate token
	t, err := middleware.GenerateTokenString(int(userPtr.ID), userPtr.Email)
	if err != nil {
		return nil, errors.New("couldn't generate token")
	}

	// send token as response
	return &t, nil
}
