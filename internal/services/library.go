package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func AddToLibrary(ctx context.Context, uid uint, bid uint, state database.ReadState) error {
	return repositories.AddBook(ctx, uid, bid, state)
}

func QueryLibrary(ctx context.Context, uid uint, limit uint64) ([]dtos.LibraryResponse, error) {
	libs, err := repositories.QueryLibraryDb(ctx, uid, limit)
	if err != nil {
		return []dtos.LibraryResponse{}, err
	}

	return convertLibtoResp(libs)
}
func convertLibtoResp(p []database.Library) ([]dtos.LibraryResponse, error) {
	res := make([]dtos.LibraryResponse, len(p))
	for i, lib := range p {
		res[i] = dtos.LibraryResponse{
			Uid:    lib.UserID,
			Book:   lib.Book,
			State:  lib.State,
			Rating: lib.Rating,
		}
	}

	return res, nil
}
