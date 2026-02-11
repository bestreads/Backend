package database

type FollowRel struct {
	User User
	// der user, der folgt
	UserID    uint `gorm:"not null;constraint:OnUpdate:CASCATE,OnDelete:CASCATE;uniqueIndex:idx_follow"`
	Following User
	// der gefolgte user
	FollowingID uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_follow"`
}
