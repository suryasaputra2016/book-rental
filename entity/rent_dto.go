package entity

type ShowRentsResponse struct {
	Message     string       `json:"message" validate:"required, message"`
	RentedBooks []RentedBook `json:"rented_books" validate:"required, rented_books"`
}
