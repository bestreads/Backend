package database

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	User      User           `gorm:"constraint:OnUpdate:CASCADE"`
	UserID    uint           `gorm:"uniqueIndex:idx_post"`
	Book      Book           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID    uint           `gorm:"uniqueIndex:idx_post"`
	Content   string
	State     ReadState `gorm:"-"`
	Rating    uint      `gorm:"-"`
}
