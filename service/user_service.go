package service

import (
	"errors"

	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/repo"
	"github.com/suryasaputra2016/book-rental/utils"
)

type UserService interface {
	CreateUser(*entity.CreateUserRequest) (*entity.User, error)
}

// user repository implementation with database connection
type userService struct {
	ur repo.UserRepo
}

func (us *userService) CreateUser(userPtr *entity.CreateUserRequest) (*entity.User, error) {
	// email string validation
	if !utils.IsEmailStringValid(userPtr.Email) {
		return nil, errors.New("email is not well formatted")
	}

	// password string validation
	err := utils.IsPasswordGood(userPtr.Password)
	if err != nil {
		return nil, err
	}

	// check email if it already exists
	_, err = us.ur.FindUserByEmail(userPtr.Email)
	if err != nil {
		return nil, errors.New("email is already in use")
	}

	// hash password
	passwordHash, err := utils.GenerateHash(userPtr.Password)
	if err != nil {
		return nil, errors.New("couldn't hash password")
	}

	//define new user
	newUserPtr := &entity.User{
		Email:         userPtr.Email,
		PasswordHash:  passwordHash,
		DepositAmount: 0.0,
	}

	// add new user to database and return
	newUserPtr, err = us.ur.AddUser(newUserPtr)
	if err != nil {
		return nil, err
	}
	return newUserPtr, err
}
