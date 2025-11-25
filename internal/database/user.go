package database

// AAAA wer hat sich das mit der großschreibung ausgedacht
// "User", weil gorm das automatisch in "users" umbenennt
type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"not null;unique;index"`
	Password_hash string
}
