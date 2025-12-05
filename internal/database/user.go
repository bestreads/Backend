package database

// AAAA wer hat sich das mit der großschreibung ausgedacht
// "User", weil gorm das automatisch in "users" umbenennt
type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"not null;unique;index"`
	Username      string
	Password_hash string
	Pfp           string // einfach der name der datei, ein sha256-hash? siehe https://stackoverflow.com/questions/3304588/store-user-profile-pictures-on-disk-or-in-the-database
}
