package entity

// rent book request dto
type RentBookRequest struct {
	ISBN   string `json:"isbn" validate:"required, isbn"`
	Title  string `json:"title" validate:"required, title"`
	Author string `json:"author" validate:"required, author"`
}

// rent book response dto
type RentBookResponse struct {
	Message    string           `json:"message" validate:"required, message"`
	UserData   UserResponseData `json:"user_data" validate:"required, user_data"`
	RentedBook RentedBook       `json:"rented_books" validate:"required, rented_books"`
}

// book rental data
type RentedBook struct {
	ISBN       string `json:"isbn" validate:"required, isbn"`
	Title      string `json:"title" validate:"required, title"`
	Author     string `json:"author" validate:"required, author"`
	Category   string `json:"category" validate:"required, category"`
	CopyNumber int    `json:"copy_number" validate:"required, copy_number"`
}
