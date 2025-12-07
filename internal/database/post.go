package database

type Post struct {
	User    User `gorm:"constraint:OnUpdate:CASCADE"`
	UserID  uint
	Book    Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BookID  uint
	Content string
	Image   string
}
