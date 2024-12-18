package entity

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, password"`
}

type CreateUserRepsonse struct {
	Message     string                 `json:"message" validate:"required, message"`
	NewUserData CreateUserResponseData `json:"new_user_data" validate:"required, new_user_data"`
}

type CreateUserResponseData struct {
	Email         string  `json:"email" validate:"required, email"`
	DepositAmount float32 `json:"deposit_amount" validate:"required, deposit_amount"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, password"`
}
