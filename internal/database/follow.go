package database

type FollowRel struct {
	User        User
	UserID      uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_follow`
	Following   User
	FollowingID uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_follow`
}
