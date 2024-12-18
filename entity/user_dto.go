package entity

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required, first_name"`
	LastName  string `json:"last_name" validate:"required, last_name"`
	Email     string `json:"email" validate:"required, email"`
	Password  string `json:"password" validate:"required, password"`
}

type CreateUserRepsonse struct {
	Message  string           `json:"message" validate:"required, message"`
	UserData UserResponseData `json:"user_data" validate:"required, user_data"`
}

type UserResponseData struct {
	FirstName     string  `json:"first_name" validate:"required, first_name"`
	LastName      string  `json:"last_name" validate:"required, last_name"`
	Email         string  `json:"email" validate:"required, email"`
	DepositAmount float32 `json:"deposit_amount" validate:"required, deposit_amount"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, password"`
}

type LoginResponse struct {
	Message string `json:"message" validate:"required, message"`
	Token   string `json:"token" validate:"required, token"`
}

type TopUpRequest struct {
	TopupAmount float32 `json:"topup_amount" validate:"required, topup_amount"`
}

type TopUpResponse CreateUserRepsonse
