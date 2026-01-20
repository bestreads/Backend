package dtos

import "github.com/bestreads/Backend/internal/database"

type PostResponse struct {
	ProfilePicture string
	Username       string
	Uid            uint
	Book           database.Book
	Content        string
}
