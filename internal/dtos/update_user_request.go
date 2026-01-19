package dtos

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=3,max=32,regexp=^[a-zA-Z0-9._-]+$"`
	Password string `json:"password" validate:"omitempty,min=12"`
	Pfp      []byte `json:"-"`
}

// IsEmpty prüft ob alle Felder leer sind
func (u *UpdateUserRequest) IsEmpty() bool {
	return u.Email == "" && u.Username == "" && u.Password == "" && len(u.Pfp) == 0
}
