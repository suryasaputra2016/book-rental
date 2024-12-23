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
	RentedBook RentedBook       `json:"rented_book" validate:"required, rented_book"`
}

// rented book data
type RentedBook struct {
	Title        string `json:"title" validate:"required, title"`
	Author       string `json:"author" validate:"required, author"`
	CopyNumber   int    `json:"copy_number" validate:"required, copy_number"`
	CheckoutDate string `json:"checkout_date" validate:"required, checkout_date"`
	DueDate      string `json:"due_date" validate:"required, due_date"`
	EndDate      string `json:"end_date" validate:"required, end_date"`
	RentStatus   string `json:"rent_status" validate:"required, rent_status"`
}

// return book response dto
type ReturnBookResponse struct {
	Message      string       `json:"message" validate:"required, message"`
	ReturnedBook ReturnedBook `json:"returned_book" validate:"required, returned_book"`
}

// return book request dto
type ReturnBookRequest struct {
	ISBN       string `json:"isbn" validate:"required, isbn"`
	Title      string `json:"title" validate:"required, title"`
	Author     string `json:"author" validate:"required, author"`
	CopyNumber int    `json:"copy_number" validate:"required, copy_number"`
}

// returned book data
type ReturnedBook struct {
	Title      string `json:"title" validate:"required, title"`
	Author     string `json:"author" validate:"required, author"`
	CopyNumber int    `json:"copy_number" validate:"required, copy_number"`
	RentStatus string `json:"rent_status" validate:"required, rent_status"`
}

// show book response dto
type ShowBooksResponse struct {
	ISBN       string  `json:"isbn" validate:"required, isbn"`
	Title      string  `json:"title" validate:"required, title"`
	Author     string  `json:"author" validate:"required, author"`
	CopyNumber int     `json:"copy_number" validate:"required, copy_number"`
	Status     string  `json:"status" validate:"required, status"`
	Category   string  `json:"category" validate:"required, category"`
	RentalCost float32 `json:"rental_cost" validate:"required, rental_cost"`
}
