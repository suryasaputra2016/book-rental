package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
)

type RentService interface {
	ShowRents(echo.Context) error
}

// user repository implementation with database connection
type rentService struct {
	rr repo.RentRepo
}

func NewRentService(rr repo.RentRepo) *rentService {
	return &rentService{
		rr: rr,
	}
}

func (rs *rentService) ShowRents(c echo.Context) error {
	// get user together with book and copy
	userId := middlewares.GetUserID(c.Get("user"))
	RentsPtr, err := rs.rr.FindRentsByUserID(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "user cannot be found")
	}

	// define res
	var rentedBooks []entity.RentedBook
	for _, rent := range *RentsPtr {
		rentedBook := entity.RentedBook{
			ISBN:       rent.BookCopy.Book.ISBN,
			Title:      rent.BookCopy.Book.Title,
			Author:     rent.BookCopy.Book.Author,
			Category:   rent.BookCopy.Book.Category,
			CopyNumber: rent.BookCopy.CopyNumber,
		}
		rentedBooks = append(rentedBooks, rentedBook)
	}

	res := entity.ShowRentsResponse{
		Message:     "Books currently being rented",
		RentedBooks: rentedBooks,
	}

	return c.JSON(http.StatusOK, res)
}
