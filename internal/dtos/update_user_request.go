package dtos

import "fmt"

type UpdateUserRequest struct {
	Email          string `json:"email" validate:"omitempty,email"`
	Username       string `json:"username" validate:"omitempty,min=3,max=32"`
	Password       string `json:"password" validate:"omitempty,min=12"`
	ProfilePicture []byte `json:"-"`
}

// IsEmpty prüft ob alle Felder leer sind
func (u *UpdateUserRequest) IsEmpty() bool {
	fmt.Printf("asd: %d", len(u.ProfilePicture))
	return u.Email == "" && u.Username == "" && u.Password == "" && len(u.ProfilePicture) == 0
}
