package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/suryasaputra2016/book-rental/entity"
	"github.com/suryasaputra2016/book-rental/repo"
)

// user service interface
type BookService interface {
	CheckBookRentalRequirements(title, author string, userID int) (*entity.Book, *entity.User, error)
	ProcessBookRental(bookPtr *entity.Book, userPtr *entity.User) (*entity.Rent, error)
	CheckBookReturnRequirements(title, author string, copyNumber, userID int) (*entity.Rent, error)
	ProcessBookReturn(rentPtr *entity.Rent) (*entity.BookCopy, error)
	GetAllBooks() (*[]entity.BookCopy, error)
	StoreRentHistory(userID uint, bookCopyID uint) error
}

// user service implementation with use repo
type bookService struct {
	ur repo.UserRepo
	rr repo.RentRepo
	br repo.BookRepo
}

// NewUserService takes user repo and gives new user service
func NewBookService(ur repo.UserRepo, rr repo.RentRepo, br repo.BookRepo) *bookService {
	return &bookService{ur: ur, rr: rr, br: br}
}

// CheckBookRentalRequirements verifies book titles, author, and user id, returns book and user upon validation
func (bs *bookService) CheckBookRentalRequirements(title, author string, userID int) (*entity.Book, *entity.User, error) {
	// find book
	bookPtr, err := bs.br.FindAvailableBookByTitleAuthor(title, author)
	if err != nil {
		return nil, nil, fmt.Errorf("checking book rental requirements: %w", err)
	}

	// get user
	userPtr, err := bs.ur.FindUserByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("checking book rental requirements: %w", err)
	}

	// check deposit amount
	if userPtr.DepositAmount < bookPtr.RentalCost {
		return nil, nil, errors.New("checking book rental requirements: insufficient deposit")
	}

	return bookPtr, userPtr, nil
}

// ProcessBookRental finishes book rental process
func (bs *bookService) ProcessBookRental(bookPtr *entity.Book, userPtr *entity.User) (*entity.Rent, error) {
	// flag the copy as rented
	rentedCopyPtr := &bookPtr.BookCopies[0]
	rentedCopyPtr.Status = "rented"
	if err := bs.br.EditBookCopy(rentedCopyPtr); err != nil {
		return nil, fmt.Errorf("processing book rental: %w", err)
	}

	// subtract user deposit and add book copy to user
	userPtr.DepositAmount -= bookPtr.RentalCost

	//create new rent and append it to user
	newRent := entity.Rent{
		UserID:   userPtr.ID,
		Status:   "ongoing",
		DueDate:  time.Now().AddDate(0, 0, 14),
		BookCopy: *rentedCopyPtr,
	}
	userPtr.Rents = append(userPtr.Rents, newRent)
	if err := bs.ur.EditUser(userPtr); err != nil {
		return nil, fmt.Errorf("processing book rental: %w", err)
	}
	return &newRent, nil
}

// get all books returns all books from database
func (bs *bookService) GetAllBooks() (*[]entity.BookCopy, error) {
	var bookCopiesPtr *[]entity.BookCopy
	bookCopiesPtr, err := bs.br.FindAllBook()
	if err != nil {
		return nil, fmt.Errorf("getting all books: %w", err)
	}
	return bookCopiesPtr, nil
}

// CheckBookReturnRequirements verifies book titles, author, copy number and user id, returns new rent upon validation
func (bs *bookService) CheckBookReturnRequirements(title, author string, copyNumber, userID int) (*entity.Rent, error) {
	// find user's rents
	rentsPtr, err := bs.rr.FindRentsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("checking book return requirements: %w", err)
	}

	fmt.Printf("\n\n\n%v\n\n\n", *rentsPtr)

	// check if the copy is rented
	var i int
	found := false
	for i = 0; i < len(*rentsPtr); i++ {
		if (*rentsPtr)[i].BookCopy.Book.Title == title &&
			(*rentsPtr)[i].BookCopy.Book.Author == author &&
			(*rentsPtr)[i].BookCopy.CopyNumber == copyNumber &&
			((*rentsPtr)[i].Status == "ongoing" || (*rentsPtr)[i].BookCopy.Status == "overdue") {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("checking book return requirements: copy not found in rents")
	}

	return &(*rentsPtr)[i], nil
}

// ProcessBookReturn finishes book return process
func (bs *bookService) ProcessBookReturn(rentPtr *entity.Rent) (*entity.BookCopy, error) {
	// change rent status
	rentPtr.Status = "finished"
	if err := bs.rr.EditRent(rentPtr); err != nil {
		return nil, fmt.Errorf("processing book return: %w", err)
	}

	// flag the copy as available
	copyPtr := &rentPtr.BookCopy
	copyPtr.Status = "available"
	err := bs.br.EditBookCopy(copyPtr)
	if err != nil {
		return nil, fmt.Errorf("processing book return: %w", err)
	}

	return copyPtr, nil
}

// StoreRentHistory updates rental history
func (bs *bookService) StoreRentHistory(userID uint, bookCopyID uint) error {
	// create rental history
	newRentalHistory := entity.RentalHistory{
		UserID:     userID,
		BookCopyID: bookCopyID,
		Type:       "take",
	}

	// store in databse
	err := bs.rr.AddRentHistory(&newRentalHistory)
	if err != nil {
		return fmt.Errorf("saving rent history: %w", err)
	}
	return nil
}
