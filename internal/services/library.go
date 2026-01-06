package services

import (
	"context"
	"fmt"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

func AddToLibrary(ctx context.Context, uid uint, bid uint, state database.ReadState) error {
	if state < 0 || state > 2 {
		return fmt.Errorf("test")
	}

	return repositories.AddBook(ctx, uid, bid, state)
}

func QueryLibrary(ctx context.Context, uid uint, limit int64) ([]dtos.LibraryResponse, error) {
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

func UpdateReadState(ctx context.Context, uid uint, bid uint, state database.ReadState) error {
	c, e := repositories.UpdateReadState(ctx, uid, bid, state)
	if c > 1 {
		fmt.Printf("uid: %d, bid: %d, state: %v, rows updated(!) > 1: %d\n", uid, bid, state, c)
		panic("THIS SHOULD NEVER HAPPEN")
	}

	return e
}

func DeleteFromLibrary(ctx context.Context, uid uint, bid uint) error {
	c, e := repositories.DeleteFromLibrary(ctx, uid, bid)
	if c > 1 {
		fmt.Printf("uid: %d, bid: %d, rows updated(!) > 1: %d\n", uid, bid, c)
		panic("THIS SHOULD NEVER HAPPEN")
	}

	return e
}
