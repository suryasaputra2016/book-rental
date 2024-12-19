package repo

import (
	"github.com/suryasaputra2016/book-rental/entity"
	"gorm.io/gorm"
)

// book repository interface
type RentRepo interface {
	AddRentHistory(*entity.RentalHistory) error
	FindRentsByUserID(userID int) (*[]entity.Rent, error)
	EditRents(*[]entity.Rent) (*[]entity.Rent, error)
	EditRent(*entity.Rent) (*entity.Rent, error)
}

// book repository implementation with database connection
type rentRepo struct {
	db *gorm.DB
}

// NewBookRepo returns new book repository implementation
func NewRentRepo(db *gorm.DB) *rentRepo {
	return &rentRepo{db: db}
}

func (rr *rentRepo) AddRentHistory(rentPtr *entity.RentalHistory) error {
	result := rr.db.Create(rentPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (rr *rentRepo) FindRentsByUserID(userID int) (*[]entity.Rent, error) {
	var rentPtr = new([]entity.Rent)
	result := rr.db.Preload("BookCopy.Book").Where("user_id = ?", userID).Find(rentPtr)
	if result.Error != nil {
		return nil, result.Error
	}
	return rentPtr, nil
}

func (rr *rentRepo) EditRents(rentPtr *[]entity.Rent) (*[]entity.Rent, error) {
	result := rr.db.Save(rentPtr)
	if result.Error != nil {
		return nil, result.Error
	}
	return rentPtr, nil
}

func (rr *rentRepo) EditRent(rentPtr *entity.Rent) (*entity.Rent, error) {
	result := rr.db.Save(rentPtr)
	if result.Error != nil {
		return nil, result.Error
	}
	return rentPtr, nil
}
