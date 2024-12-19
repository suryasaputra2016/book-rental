package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
)

type BookHandler interface {
	RentBook(echo.Context) error
	ReturnBook(echo.Context) error
	ShowBooks(echo.Context) error
}

// user repository implementation with database connection
type bookHandler struct {
	br repo.BookRepo
	ur repo.UserRepo
	rr repo.RentRepo
}

func NewBookHandler(br repo.BookRepo, ur repo.UserRepo, rr repo.RentRepo) *bookHandler {
	return &bookHandler{
		br: br,
		ur: ur,
		rr: rr,
	}
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
func (bs *bookHandler) RentBook(c echo.Context) error {
	// bind request body
	var req entity.RentBookRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// find available book
	var bookPtr *entity.Book
	var err error
	if bookPtr, err = bs.br.FindAvailableBookByTitleAuthor(req.Title, req.Author); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprint(err))
	}

	// get user
	userID := middlewares.GetUserID(c.Get("user"))
	var userPtr *entity.User
	if userPtr, err = bs.ur.FindUserByID(userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "user id cannot be found")
	}

	// check deposit amount
	if userPtr.DepositAmount < bookPtr.RentalCost {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient deposit, please top up")
	}

	// flag the copy as rented
	rentedCopyPtr := &bookPtr.BookCopies[0]
	rentedCopyPtr.Status = "rented"
	if rentedCopyPtr, err = bs.br.EditBookCopy(rentedCopyPtr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update book copy")
	}

	// subtract user deposit and add book copy to user
	userPtr.DepositAmount -= bookPtr.RentalCost

	//create new rent and append it to user
	newRent := entity.Rent{
		UserID:   uint(userID),
		Status:   "ongoing",
		DueDate:  time.Now().AddDate(0, 0, 14),
		BookCopy: *rentedCopyPtr,
	}
	userPtr.Rents = append(userPtr.Rents, newRent)
	if userPtr, err = bs.ur.EditUser(userPtr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update user")
	}

	// record in rental history
	newRentalHistory := entity.RentalHistory{
		UserID:     uint(userID),
		BookCopyID: *&rentedCopyPtr.BookID,
		Type:       "take",
	}
	if err = bs.rr.AddRentHistory(&newRentalHistory); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update rental history")
	}

	// define and send response
	res := entity.RentBookResponse{
		Message: "book is successfully rented",
		UserData: entity.UserResponseData{
			FirstName:     userPtr.FirstName,
			LastName:      userPtr.LastName,
			Email:         userPtr.Email,
			DepositAmount: userPtr.DepositAmount,
		},
		RentedBook: entity.RentedBook{
			Title:        bookPtr.Title,
			Author:       bookPtr.Author,
			CopyNumber:   rentedCopyPtr.CopyNumber,
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
// @Router /books/return [post]
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object}  entity.ErrorMessage
func (bs *bookHandler) ReturnBook(c echo.Context) error {
	// bind request body
	var req entity.ReturnBookRequest
	if c.Bind(&req) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "JSON request is invalid")
	}

	// find available book
	userID := middlewares.GetUserID(c.Get("user"))
	rentsPtr, err := bs.rr.FindRentsByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "rents cannot be found")
	}

	// check if the copy is rented
	var i int
	found := false
	for i = 0; i < len(*rentsPtr); i++ {
		if (*rentsPtr)[i].BookCopy.Book.Title == req.Title &&
			(*rentsPtr)[i].BookCopy.Book.Author == req.Author &&
			(*rentsPtr)[i].BookCopy.CopyNumber == req.CopyNumber &&
			((*rentsPtr)[i].BookCopy.Status == "ongoing" || (*rentsPtr)[i].BookCopy.Status == "overdue") {
			found = true
			break
		}
	}
	if !found {
		return echo.NewHTTPError(http.StatusBadRequest, "book is currently not rented by user")
	}

	// change rent status
	(*rentsPtr)[i].Status = "finished"
	if _, err = bs.rr.EditRent(&(*rentsPtr)[i]); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update rent")
	}

	// flag the copy as available
	rentedCopyPtr := &(*rentsPtr)[i].BookCopy
	rentedCopyPtr.Status = "available"
	if rentedCopyPtr, err = bs.br.EditBookCopy(rentedCopyPtr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update book copy")
	}

	// record in rental history
	newRentalHistory := entity.RentalHistory{
		UserID:     uint(userID),
		BookCopyID: *&rentedCopyPtr.BookID,
		Type:       "return",
	}
	if err = bs.rr.AddRentHistory(&newRentalHistory); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update rental history")
	}

	// define and send response
	res := entity.ReturnBookResponse{
		Message: "book is successfully returned",
		ReturnedBook: entity.ReturnedBook{
			Title:      rentedCopyPtr.Book.Title,
			Author:     rentedCopyPtr.Book.Author,
			CopyNumber: rentedCopyPtr.CopyNumber,
			RentStatus: rentedCopyPtr.Status,
		},
	}
	return c.JSON(http.StatusOK, res)
}

// @Summary Show All Book
// @Description Show All Book in the library
// @Tags books
// @Produce json
// @Security JWT
// @Success 200 {object} entity.ShowBooksResponse
// @Router /books [get]
// @Failure 500 {object}  entity.ErrorMessage
func (bs *bookHandler) ShowBooks(c echo.Context) error {
	var bookCopies []entity.BookCopy
	bookCopies, err := bs.br.FindAllBook()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get books")
	}

	var res []entity.ShowBooksResponse
	var copyResponse entity.ShowBooksResponse
	for _, copy := range bookCopies {
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
