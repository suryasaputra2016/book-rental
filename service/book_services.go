package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/middlewares"
	"github.com/suryasaputra2016/book-rental/repo"
)

type BookService interface {
	RentBook(echo.Context) error
}

// user repository implementation with database connection
type bookService struct {
	br repo.BookRepo
	ur repo.UserRepo
	rr repo.RentRepo
}

func NewBookService(br repo.BookRepo, ur repo.UserRepo, rr repo.RentRepo) *bookService {
	return &bookService{
		br: br,
		ur: ur,
		rr: rr,
	}
}

func (bs *bookService) RentBook(c echo.Context) error {
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
		DueDate:  time.Now(),
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
			ISBN:       bookPtr.ISBN,
			Title:      bookPtr.Title,
			Author:     bookPtr.Author,
			Category:   bookPtr.Category,
			CopyNumber: rentedCopyPtr.CopyNumber,
		},
	}
	return c.JSON(http.StatusOK, res)
}
