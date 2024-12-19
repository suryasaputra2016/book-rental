package services

import (
	"errors"
	"fmt"

	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/utils"
)

// user service interface
type UserService interface {
	CheckRegistrationData(email, password string) error
	CreateNewUser(registrationData *entity.RegisterRequest) (*entity.User, error)
	CheckLoginData(email, password string) (*entity.User, error)
	UpdateDeposit(userPtr *entity.User, amount float32) error
	CheckTopupData(userID int, topupAmount float32) (*entity.User, error)
}

// user service implementation with use repo
type userService struct {
	ur repo.UserRepo
}

// NewUserService takes user repo and gives new user service
func NewUserService(ur repo.UserRepo) *userService {
	return &userService{ur: ur}
}

// CheckRegistrationData checks if email, and password are well formatted and email hasn't been used
func (us *userService) CheckRegistrationData(email, password string) error {
	// email validation
	if err := utils.IsEmailStringValid(email); err != nil {
		return fmt.Errorf("validating registration data: %w", err)
	}

	// password validation
	if err := utils.IsPasswordGood(password); err != nil {
		return fmt.Errorf("validating registration data: %w", err)
	}

	// check email if it already exists
	if _, err := us.ur.FindUserByEmail(email); err == nil {
		return fmt.Errorf("validating registration data: %w", err)
	}

	return nil
}

// CreateNewUser accepts registration data and returns new user
func (us *userService) CreateNewUser(registrationData *entity.RegisterRequest) (*entity.User, error) {
	// hash password
	passwordHash, err := utils.GenerateHash(registrationData.Password)
	if err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	// define new user
	newUser := entity.User{
		FirstName:     registrationData.FirstName,
		LastName:      registrationData.LastName,
		Email:         registrationData.Email,
		PasswordHash:  passwordHash,
		DepositAmount: 0.0,
	}

	// add new user to database
	if err := us.ur.AddUser(&newUser); err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	return &newUser, nil
}

// CheckLoginData verifies email, and password to database and if verified returns user
func (us *userService) CheckLoginData(email, password string) (*entity.User, error) {
	// check email and get user
	userPtr, err := us.ur.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("validating login data: %w", err)
	}

	// check password
	if err := utils.CheckPassword(password, userPtr.PasswordHash); err != nil {
		return nil, fmt.Errorf("validating login data: %w", err)
	}

	return userPtr, nil
}

// UpdateDeposit update user's deposit amount and save it on database
func (us *userService) UpdateDeposit(userPtr *entity.User, amount float32) error {
	// update user deposit amount
	userPtr.DepositAmount += amount

	// save to database
	if err := us.ur.EditUser(userPtr); err != nil {
		return fmt.Errorf("updating deposit amount: %w", err)
	}

	return nil
}

// CheckTopupData verifies user id and top-up amount and if verified returns user
func (us *userService) CheckTopupData(userID int, topupAmount float32) (*entity.User, error) {
	// check user id and get user
	userPtr, err := us.ur.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("validating top-up data: %w", err)
	}

	//validate req.Amount
	if topupAmount <= 0.0 {
		return nil, errors.New("validating top-up data: non-positif top-up amoung")
	}

	return userPtr, nil
}
