package dtos

import "github.com/bestreads/Backend/internal/database"

type ActivityResponse struct {
	Posts    []PostResponse
	Activity []database.ReadState
}
