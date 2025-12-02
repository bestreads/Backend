package database

type state int

const (
	Unread state = iota
	Reading
	Read
)

type RelBookUser struct {
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID uint
	Book   Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID uint
	State  state
	Rating uint
}
