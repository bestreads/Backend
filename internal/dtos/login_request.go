package dtos

type LoginRequest struct {
	Email    string `json:"email" example:"test@bestreads.byte-flow.de" validate:"required,email"`
	Password string `json:"password" example:"ichBinEinPasswort" validate:"required,min=12"`
}
