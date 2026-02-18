package database

import (
	"gorm.io/gorm"
)

// AAAA wer hat sich das mit der großschreibung ausgedacht
// "User", weil gorm das automatisch in "users" umbenennt
type User struct {
	gorm.Model
	Email          string `gorm:"not null;unique;index"`
	Password_hash  string
	Username       string `gorm:"not null"`
	ProfilePicture string
	Description    string
}
