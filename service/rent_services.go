package service

import (
	"net/http"
	"time"

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

// @Summary Show Rents
// @Description Show all rents, finished and ongoing
// @Tags rents
// @Produce json
// @Security JWT
// @Success 200 {object} entity.ShowRentsResponse
// @Router /rents [get]
// @Failure 500 {object}  entity.ErrorMessage
func (rs *rentService) ShowRents(c echo.Context) error {
	// get rents and copy
	userId := middlewares.GetUserID(c.Get("user"))
	RentsPtr, err := rs.rr.FindRentsByUserID(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "rents cannot be found")
	}

	// check rent status and update
	var i int
	for i = 0; i < len(*RentsPtr); i++ {
		if (*RentsPtr)[i].DueDate.Before(time.Now()) {
			(*RentsPtr)[i].Status = "overdue"
		}
	}
	if i > 0 {
		if RentsPtr, err = rs.rr.EditRents(RentsPtr); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "cannot update rents")
		}
	}

	// define res
	var rentedBooks []entity.RentedBook
	for _, rent := range *RentsPtr {
		rentedBook := entity.RentedBook{
			Title:        rent.BookCopy.Book.Title,
			Author:       rent.BookCopy.Book.Author,
			CopyNumber:   rent.BookCopy.CopyNumber,
			CheckoutDate: rent.StartDate.Format("2006-01-02"),
			DueDate:      rent.DueDate.Format("2006-01-02"),
			RentStatus:   rent.Status,
		}
		rentedBooks = append(rentedBooks, rentedBook)
	}

	res := entity.ShowRentsResponse{
		Message:     "Books currently being rented",
		RentedBooks: rentedBooks,
	}

	return c.JSON(http.StatusOK, res)
}
