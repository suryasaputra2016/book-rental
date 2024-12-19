package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/utils"
)

type UserService interface {
	CreateUser(echo.Context) error
	Login(echo.Context) error
	Topup(echo.Context) error
	ShowRents(echo.Context) error
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
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		DepositAmount: 0.0,
	}

	// add new user to database and return
	if us.ur.AddUser(&newUser) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	//define and send response
	res := entity.CreateUserRepsonse{
		Message: "user successfully created",
		UserData: entity.UserResponseData{
			FirstName:     newUser.FirstName,
			LastName:      newUser.LastName,
			Email:         newUser.Email,
			DepositAmount: newUser.DepositAmount,
		},
	}
	return c.JSON(http.StatusCreated, res)
}

func (us *userService) Login(c echo.Context) error {
	// bind request body
	var req entity.LoginRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// check email and get user info
	var userPtr *entity.User
	var err error
	if userPtr, err = us.ur.FindUserByEmail(req.Email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "email cannot be found")
	}

	// check password
	if utils.CheckPassword(req.Password, userPtr.PasswordHash) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "password doesn't match")
	}

	// generate token
	t, err := middlewares.GenerateTokenString(int(userPtr.ID), userPtr.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't generate token")
	}

	// define and send response
	res := entity.LoginResponse{
		Message: "login successful",
		Token:   t,
	}
	return c.JSON(http.StatusOK, res)
}

func (us *userService) Topup(c echo.Context) error {
	// get res
	// bind request body
	var req entity.TopUpRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// get user id
	userId := middlewares.GetUserID(c.Get("user"))

	// find user with id
	var userPtr *entity.User
	var err error
	if userPtr, err = us.ur.FindUserByID(userId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "id cannot be found")
	}

	// xendit

	//validate req.Amount
	if req.TopupAmount <= 0.0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "top up amount must be positive float")
	}

	// update user deposit amount
	userPtr.DepositAmount += req.TopupAmount

	// save to database
	if userPtr, err = us.ur.EditUser(userPtr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update user")
	}

	// define and send response
	res := entity.TopUpResponse{
		Message: "deposit amount is successfully updated",
		UserData: entity.UserResponseData{
			FirstName:     userPtr.FirstName,
			LastName:      userPtr.LastName,
			Email:         userPtr.Email,
			DepositAmount: userPtr.DepositAmount,
		},
	}
	return c.JSON(http.StatusOK, res)
}
