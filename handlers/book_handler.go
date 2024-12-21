package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/services"
)

// book handler interface
type BookHandler interface {
	RentABook(echo.Context) error
	ReturnABook(echo.Context) error
	ShowBooks(echo.Context) error
}

// boook handler implementation with book service
type bookHandler struct {
	bs services.BookService
}

// NewBookHandler takes book service and returns new book handler
func NewBookHandler(bs services.BookService) *bookHandler {
	return &bookHandler{bs: bs}
}

// @Summary Rent A Book
// @Description Rent One Book
// @Tags books
// @Accept json
// @Produce json
// @Param rent-book-data body entity.RentBookRequest true "rentbook request"
// @Security JWT
// @Success 200 {object} entity.RentBookResponse
// @Router /books/rent [post]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (bh *bookHandler) RentABook(c echo.Context) error {
	// bind request body
	var req entity.RentBookRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON request")
	}

	// check book rental requirements
	userID := middlewares.GetUserID(c.Get("user"))
	bookPtr, userPtr, err := bh.bs.CheckBookRentalRequirements(req.Title, req.Author, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// process book rental
	newRent, err := bh.bs.ProcessBookRental(bookPtr, userPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// record in rent history
	if err = bh.bs.StoreRentHistory(uint(userID), bookPtr.BookCopies[0].ID, "take"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// define and send response
	res := entity.RentBookResponse{
		Message: "rental success",
		UserData: entity.UserResponseData{
			FirstName:     userPtr.FirstName,
			LastName:      userPtr.LastName,
			Email:         userPtr.Email,
			DepositAmount: userPtr.DepositAmount,
		},
		RentedBook: entity.RentedBook{
			Title:        bookPtr.Title,
			Author:       bookPtr.Author,
			CopyNumber:   bookPtr.BookCopies[0].CopyNumber,
			CheckoutDate: time.Now().Format("2006-01-02"),
			DueDate:      newRent.DueDate.Format("2006-01-02"),
			RentStatus:   newRent.Status,
		},
	}
	return c.JSON(http.StatusOK, res)
}

// @Summary Return A Book
// @Description Return One Book
// @Tags books
// @Accept json
// @Produce json
// @Param return-book-data body entity.ReturnBookRequest true "returnbook request"
// @Security JWT
// @Success 200 {object} entity.ReturnBookResponse
// @Router /books/return [put]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (bh *bookHandler) ReturnABook(c echo.Context) error {
	// bind request body
	var req entity.ReturnBookRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// check book return requirements
	userID := middlewares.GetUserID(c.Get("user"))
	rentPtr, err := bh.bs.CheckBookReturnRequirements(req.Title, req.Author, req.CopyNumber, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// process book return
	copyPtr, err := bh.bs.ProcessBookReturn(rentPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// record in rent history
	if err = bh.bs.StoreRentHistory(uint(userID), copyPtr.ID, "return"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// define and send response
	res := entity.ReturnBookResponse{
		Message: "return success",
		ReturnedBook: entity.ReturnedBook{
			Title:      copyPtr.Book.Title,
			Author:     copyPtr.Book.Author,
			CopyNumber: copyPtr.CopyNumber,
			RentStatus: copyPtr.Status,
		},
	}
	return c.JSON(http.StatusOK, res)
}

// @Summary Show All Books
// @Description Show All Books in the library
// @Tags books
// @Produce json
// @Success 200 {object} entity.ShowBooksResponse
// @Router /books [get]
// @Failure 500 {object}  entity.ErrorMessage
func (bh *bookHandler) ShowBooks(c echo.Context) error {
	// get all books
	bookCopiesPtr, err := bh.bs.GetAllBooks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// define and send response
	var res []entity.ShowBooksResponse
	var copyResponse entity.ShowBooksResponse
	for _, copy := range *bookCopiesPtr {
		copyResponse = entity.ShowBooksResponse{
			ISBN:       copy.Book.ISBN,
			Title:      copy.Book.Title,
			Author:     copy.Book.Author,
			Category:   copy.Book.Category,
			RentalCost: copy.Book.RentalCost,
			CopyNumber: copy.CopyNumber,
			Status:     copy.Status,
		}
		res = append(res, copyResponse)
	}
	return c.JSON(http.StatusOK, res)
}
