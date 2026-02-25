package dtos

import (
	"time"

	"github.com/bestreads/Backend/internal/database"
)

type PostResponse struct {
	ProfilePicture string
	Username       string
	Uid            uint
	BookID         uint
	Book           database.Book
	Content        string
	UpdatedAt      time.Time
	State          database.ReadState
	Rating         uint
}
