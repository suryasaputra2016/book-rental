package service

import (
	"fmt"
	"net/http"

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
}

func NewBookService(br repo.BookRepo, ur repo.UserRepo) *bookService {
	return &bookService{
		br: br,
		ur: ur,
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
	userId := middlewares.GetUserID(c.Get("user"))
	var userPtr *entity.User
	if userPtr, err = bs.ur.FindUserByID(userId); err != nil {
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
	userPtr.BookCopies = append(userPtr.BookCopies, *rentedCopyPtr)
	if userPtr, err = bs.ur.EditUser(userPtr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot update user")
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
