package repo

import (
	"github.com/suryasaputra2016/book-rental/entity"
	"gorm.io/gorm"
)

// user repository interface
type UserRepo interface {
	FindUserByID(id int) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	AddUser(userPtr *entity.User) error
	EditUser(userPtr *entity.User) error
}

// user repository implementation with database connection
type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

// FindUserByID gets user with id
func (ur *userRepo) FindUserByID(id int) (*entity.User, error) {
	var userPtr = new(entity.User)
	result := ur.db.First(userPtr, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPtr, nil
}

// FindUserByEmail gets user with email
func (ur *userRepo) FindUserByEmail(email string) (*entity.User, error) {
	var userPtr = new(entity.User)
	result := ur.db.Where("email = ?", email).First(userPtr)
	if result.Error != nil {
		return nil, result.Error
	}
	return userPtr, nil
}

// AddUser inserts user to database
func (ur *userRepo) AddUser(userPtr *entity.User) error {
	result := ur.db.Create(userPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// EditUser updates user in database
func (ur *userRepo) EditUser(userPtr *entity.User) error {
	result := ur.db.Save(userPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
