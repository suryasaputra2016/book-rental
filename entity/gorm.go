package entity

import "time"

// user entity
type User struct {
	ID            uint    `gorm:"primaryKey"`
	FirstName     string  `gorm:"type:varchar(100);not null"`
	LastName      string  `gorm:"type:varchar(100);not null"`
	Email         string  `gorm:"type:varchar(100);not null;unique"`
	PasswordHash  string  `gorm:"type:text;not null"`
	DepositAmount float32 `gorm:"type:decimal(15,2); default:0.0"`
	BookCopies    []BookCopy
}

// book entity
type Book struct {
	ID         uint    `gorm:"primaryKey"`
	ISBN       string  `gorm:"type:varchar(13);not null; unique"`
	Title      string  `gorm:"type:varchar(255);not null"`
	Author     string  `gorm:"type:varchar(255);not null"`
	Category   string  `gorm:"type:varchar(255);not null"` //comic, novel,biography, art, textbook
	RentalCost float32 `gorm:"type:decimal(15,2);not null"`
	BookCopies []BookCopy
}

// book copy entity
type BookCopy struct {
	ID         uint `gorm:"primaryKey"`
	BookID     uint
	CopyNumber int    `gorm:"not null"`
	Status     string `gorm:"type:varchar(255);not null; default:available"` //available, rented, in repair
}

// rental history
type RentalHistory struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	BookID    uint
	Type      string // rent, return
	CreatedAt time.Time
}
