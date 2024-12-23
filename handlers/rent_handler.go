package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/services"
)

// rent handler interface
type RentHandler interface {
	ShowRents(echo.Context) error
}

// rent handler implementation with rent service
type rentHandler struct {
	rs services.RentService
}

// NewRentHandler takes rent service and returns new rent handler
func NewRentHandler(rs services.RentService) *rentHandler {
	return &rentHandler{rs: rs}
}

// @Summary Show Rents
// @Description Show all rents, finished and ongoing
// @Tags rents
// @Produce json
// @Security JWT
// @Success 200 {object} entity.ShowRentsResponse
// @Router /rents [get]
// @Failure 500 {object}  entity.ErrorMessage
func (rh *rentHandler) ShowRents(c echo.Context) error {
	// get rents
	userID := middlewares.GetUserID(c.Get("user"))
	rentsPtr, err := rh.rs.GetRents(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// define and send response
	var rentedBooks []entity.RentedBook
	for _, rent := range *rentsPtr {
		var endDate string
		if rent.EndDate != nil {
			endDate = rent.EndDate.Format("2006-01-02")
		} else {
			endDate = ""
		}

		rentedBook := entity.RentedBook{
			Title:        rent.BookCopy.Book.Title,
			Author:       rent.BookCopy.Book.Author,
			CopyNumber:   rent.BookCopy.CopyNumber,
			CheckoutDate: rent.StartDate.Format("2006-01-02"),
			DueDate:      rent.DueDate.Format("2006-01-02"),
			EndDate:      endDate,
			RentStatus:   rent.Status,
		}
		rentedBooks = append(rentedBooks, rentedBook)
	}

	res := entity.ShowRentsResponse{
		Message:     "Rental Activity",
		RentedBooks: rentedBooks,
	}

	return c.JSON(http.StatusOK, res)
}
