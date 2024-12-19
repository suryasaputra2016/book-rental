package entity

type ErrorMessage struct {
	Error string `json:"error" validate:"required, error"`
}
