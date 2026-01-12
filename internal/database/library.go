package database

type ReadState int

const (
	WantToRead ReadState = iota
	Reading
	Read
)

// das ding repräsentiert einen eintrag in der persönlichen bibliothek eines nutzers
// https://gorm.io/docs/indexes.html#Composite-Indexes
type Library struct {
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID uint `gorm:"uniqueIndex:idx_lib"`
	Book   Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID uint `gorm:"uniqueIndex:idx_lib"`
	State  ReadState
	Rating uint
}
