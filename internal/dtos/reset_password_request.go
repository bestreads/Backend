package dtos

type ResetPasswordRequest struct {
	Email          string `validate:"required,email"`
	SecurityAnswer string `validate:"required,min=1"`
	NewPassword    string `validate:"required,min=12"`
}
