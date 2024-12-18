package entity

// user entity
type User struct {
	ID            uint    `gorm:"primaryKey"`
	FirstName     string  `gorm:"type:varchar(100);not null"`
	LastName      string  `gorm:"type:varchar(100);not null"`
	Email         string  `gorm:"type:varchar(100);not null;unique"`
	PasswordHash  string  `gorm:"type:text;not null"`
	DepositAmount float32 `gorm:"type:decimal(15,2); default:0.0"`
	Rents         []Rent
}

// rent entity
type Rent struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Books  []Book
}

// book entity
type Book struct {
	ID     uint `gorm:"primaryKey"`
	RentID uint
	ISBN   string `gorm:"type:varchar(13);not null"`
	Title  string `gorm:"type:varchar(255);not null"`
	Author string `gorm:"type:varchar(255);not null"`
	Copy   int    `gorm:"not null"`
	Status string `gorm:"type:varchar(255);not null; default:available"`
}
