package repo

import (
	"github.com/suryasaputra2016/book-rental/entity"
	"gorm.io/gorm"
)

// rent repository interface
type RentRepo interface {
	AddRentHistory(rentPtr *entity.RentalHistory) error
	FindRentsByUserID(userID int) (*[]entity.Rent, error)
	EditRents(rentsPtr *[]entity.Rent) error
	EditRent(rentPtr *entity.Rent) error
}

// rent repository implementation with database connection
type rentRepo struct {
	db *gorm.DB
}

// NewRentRepo takes database connection and returns new book repository implementation
func NewRentRepo(db *gorm.DB) *rentRepo {
	return &rentRepo{db: db}
}

// AddRentHistory inserts new rent history to database
func (rr *rentRepo) AddRentHistory(rentPtr *entity.RentalHistory) error {
	result := rr.db.Create(rentPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindRentsByUserID get user by id
func (rr *rentRepo) FindRentsByUserID(userID int) (*[]entity.Rent, error) {
	var rentPtr = new([]entity.Rent)
	result := rr.db.Preload("BookCopy.Book").Where("user_id = ?", userID).Find(rentPtr)
	if result.Error != nil {
		return nil, result.Error
	}
	return rentPtr, nil
}

// EditRents updates multiple rent rows on database
func (rr *rentRepo) EditRents(rentsPtr *[]entity.Rent) error {
	result := rr.db.Save(rentsPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// EditRent updates single rent row on database
func (rr *rentRepo) EditRent(rentPtr *entity.Rent) error {
	result := rr.db.Save(rentPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
