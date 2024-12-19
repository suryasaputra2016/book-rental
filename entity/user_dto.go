package entity

// reqister request dto
type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required, first_name"`
	LastName  string `json:"last_name" validate:"required, last_name"`
	Email     string `json:"email" validate:"required, email"`
	Password  string `json:"password" validate:"required, password"`
}

// reqister response dto
type RegisterRepsonse struct {
	Message  string           `json:"message" validate:"required, message"`
	UserData UserResponseData `json:"user_data" validate:"required, user_data"`
}

// user response data
type UserResponseData struct {
	FirstName     string  `json:"first_name" validate:"required, first_name"`
	LastName      string  `json:"last_name" validate:"required, last_name"`
	Email         string  `json:"email" validate:"required, email"`
	DepositAmount float32 `json:"deposit_amount" validate:"required, deposit_amount"`
}

// login request dto
type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, password"`
}

// login response dto
type LoginResponse struct {
	Message string `json:"message" validate:"required, message"`
	Token   string `json:"token" validate:"required, token"`
}

// top up request dto
type TopupRequest struct {
	TopupAmount float32 `json:"topup_amount" validate:"required, topup_amount"`
}

// top up request dto
type TopupResponse RegisterRepsonse
