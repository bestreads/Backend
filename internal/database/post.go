package database

type Post struct {
	User      UserMeta `gorm:"constraint:OnUpdate:CASCADE"`
	UserID    uint
	Book      Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID    uint
	Content   string
	ImageHash string
}
