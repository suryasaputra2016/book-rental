package entities

type User struct {
	ID            uint    `gorm:"primaryKey"`
	Email         string  `gorm:"type:varchar(100);not null;unique"`
	PasswordHash  string  `gorm:"type:text;not null"`
	DepositAmount float32 `gorm:"type:decimal(15,2); default:0.0"`
}
