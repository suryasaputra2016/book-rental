package services

import (
	"fmt"
	"time"

	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/repo"
)

// rent service interface
type RentService interface {
	GetRents(userID int) (*[]entity.Rent, error)
}

// user service implementation with use repo
type rentService struct {
	rr repo.RentRepo
}

// NewUserService takes user repo and gives new user service
func NewRentService(rr repo.RentRepo) *rentService {
	return &rentService{rr: rr}
}

// GetRents checks user id and gives updated rents related to the user
func (rs *rentService) GetRents(userID int) (*[]entity.Rent, error) {
	// check user id and get rents
	rentsPtr, err := rs.rr.FindRentsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getting rents: %w", err)
	}

	// update rents statuses
	var i int
	for i = 0; i < len(*rentsPtr); i++ {
		if (*rentsPtr)[i].DueDate.Before(time.Now()) {
			(*rentsPtr)[i].Status = "overdue"
		}
	}
	if i > 0 {
		if err = rs.rr.EditRents(rentsPtr); err != nil {
			return nil, fmt.Errorf("getting rents: %w", err)
		}
	}

	return rentsPtr, nil
}
