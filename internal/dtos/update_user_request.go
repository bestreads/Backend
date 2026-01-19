package dtos

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty,min=12"`
	Pfp      string `json:"pfp" validate:"omitempty"`
}

// IsEmpty prüft ob alle Felder leer sind
func (u *UpdateUserRequest) IsEmpty() bool {
	return u.Email == "" && u.Username == "" && u.Password == "" && u.Pfp == ""
}
