package database

// AAAA wer hat sich das mit der großschreibung ausgedacht
// "User", weil gorm das automatisch in "users" umbenennt
type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"not null;unique;index"`
	Password_hash string
	Metadata      UserMeta `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	MetadataID    uint
}

type UserMeta struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Pfp      string
}
