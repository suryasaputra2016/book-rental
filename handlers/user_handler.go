package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/services"
)

// user handler interface
type UserHandler interface {
	Register(echo.Context) error
	Login(echo.Context) error
	Topup(echo.Context) error
	ShowRents(echo.Context) error
}

// user handler implementation with user service
type userHandler struct {
	us services.UserService
}

// NewUserHandler takes user service and returns new user handler
func NewUserHandler(us services.UserService) *userHandler {
	return &userHandler{us: us}
}

// @Summary Register
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param register-data body entity.CreateUserRequest true "register request"
// @Success 201 {object} entity.CreateUserRepsonse
// @Router /register [post]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (uh *userHandler) Register(c echo.Context) error {
	// bind request body
	var req entity.RegisterRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON request")
	}

	// check registration data
	if err := uh.us.CheckRegistrationData(req.Email, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// create new user
	newUserPtr, err := uh.us.CreateNewUser(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	//define and send response
	res := entity.RegisterRepsonse{
		Message: "registration success",
		UserData: entity.UserResponseData{
			FirstName:     newUserPtr.FirstName,
			LastName:      newUserPtr.LastName,
			Email:         newUserPtr.Email,
			DepositAmount: newUserPtr.DepositAmount,
		},
	}
	return c.JSON(http.StatusCreated, res)
}

// @Summary Login
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param login-data body entity.LoginRequest true "login request"
// @Success 200 {object} entity.LoginResponse
// @Router /login [post]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (uh *userHandler) Login(c echo.Context) error {
	// bind request body
	var req entity.LoginRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON request")
	}

	// check login data
	userPtr, err := uh.us.CheckLoginData(req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// generate token
	t, err := middlewares.GenerateTokenString(int(userPtr.ID), userPtr.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't generate token")
	}

	// define and send response
	res := entity.LoginResponse{
		Message: "login success",
		Token:   t,
	}
	return c.JSON(http.StatusOK, res)
}

// @Summary Top-up
// @Description Top-up deposit amount
// @Tags user
// @Accept json
// @Produce json
// @Param topup-data body entity.TopupRequest true "topup request"
// @Security JWT
// @Success 200 {object} entity.TopupResponse
// @Router /topup [put]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (uh *userHandler) Topup(c echo.Context) error {
	// bind request body
	var req entity.TopupRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON request")
	}

	// check top up data
	userID := middlewares.GetUserID(c.Get("user"))
	userPtr, err := uh.us.CheckTopupData(userID, req.TopupAmount)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// create invoice
	// invoiceRes, err := utils.CreateInvoice(*userPtr, req.TopupAmount)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, "Error while creating invoices")
	// }
	// return c.JSON(http.StatusOK, invoiceRes)

	// update deposit
	if err := uh.us.UpdateDeposit(userPtr, req.TopupAmount); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// define and send response
	res := entity.TopupResponse{
		Message: "top up success",
		UserData: entity.UserResponseData{
			FirstName:     userPtr.FirstName,
			LastName:      userPtr.LastName,
			Email:         userPtr.Email,
			DepositAmount: userPtr.DepositAmount,
		},
	}
	return c.JSON(http.StatusOK, res)
}
