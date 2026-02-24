package dtos

type UpdateUserRequest struct {
	Email          string `json:"email" validate:"omitempty,email"`
	Username       string `json:"username" validate:"omitempty,min=3,max=32"`
	Password       string `json:"password" validate:"omitempty,min=12"`
	Description    string `json:"description" validate:"omitempty"`
	ProfilePicture []byte `json:"profile_picture"`
}

// IsEmpty prüft ob alle Felder leer sind
// flaky check
func (u *UpdateUserRequest) IsEmpty() bool {
	return u.Email == "" && u.Username == "" && u.Password == "" && len(u.ProfilePicture) == 0 && u.Description == ""
}
