package database

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	User    User `gorm:"constraint:OnUpdate:CASCADE"`
	UserID  uint
	Book    Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID  uint
	Content string
	State   ReadState `gorm:"-"`
	Rating  uint      `gorm:"-"`
}
