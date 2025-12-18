package dtos

import "github.com/bestreads/Backend/internal/database"

type LibraryResponse struct {
	Uid    uint
	Book   database.Book
	State  database.ReadState
	Rating uint
}
