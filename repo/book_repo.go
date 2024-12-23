package repo

import (
	"errors"

	"github.com/suryasaputra2016/book-rental/entity"
	"gorm.io/gorm"
)

// book repository interface
type BookRepo interface {
	FindBookByTitleAuthor(title, author string) (*entity.Book, error)
	FindAvailableBookByTitleAuthor(title, author string) (*entity.Book, error)
	FindAllBook() (*[]entity.BookCopy, error)
	EditBookCopy(bookCopyPtr *entity.BookCopy) error
}

// book repository implementation with database connection
type bookRepo struct {
	db *gorm.DB
}

// NewBookRepo takes database connection and returns new book repository implementation
func NewBookRepo(db *gorm.DB) *bookRepo {
	return &bookRepo{db: db}
}

// FindBookByTitleAuthor gets book by title and author
func (br bookRepo) FindBookByTitleAuthor(title, author string) (*entity.Book, error) {

	var bookPtr = new(entity.Book)
	result := br.db.Preload("BookCopies").
		Where("title = ? AND author = ?", title, author).
		First(bookPtr)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookPtr, nil
}

// FindAvailableBookByTitleAuthor gets available book by title and author
func (br bookRepo) FindAvailableBookByTitleAuthor(title, author string) (*entity.Book, error) {
	bookPtr, err := br.FindBookByTitleAuthor(title, author)
	if err != nil {
		return nil, errors.New("we don't have the book")
	}
	//check available copies
	var available_copies int
	for _, copy := range bookPtr.BookCopies {
		if copy.Status == "available" {
			available_copies++
		}
	}
	if available_copies == 0 {
		return nil, errors.New("no book copy is available for rent")
	}
	return bookPtr, nil
}

// FindAllBook gets all books from the database
func (br bookRepo) FindAllBook() (*[]entity.BookCopy, error) {
	var bookCopies []entity.BookCopy
	result := br.db.Preload("Book").Order("id").Find(&bookCopies)
	if result.Error != nil {
		return nil, result.Error
	}

	return &bookCopies, nil
}

// EditBookCopy updates book copy in database
func (br *bookRepo) EditBookCopy(bookCopyPtr *entity.BookCopy) error {
	result := br.db.Model(bookCopyPtr).Select("status").Save(bookCopyPtr)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
