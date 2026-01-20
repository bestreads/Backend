package services

import (
	"context"

	"github.com/bestreads/Backend/internal/repositories"
)

// RecalculateRatingAvg gets all ratings of the given book and saves the updated rating avg
func RecalculateRatingAvg(ctx context.Context, bookId uint) error {
	// Get library entries for the given book
	libraries, readLibrariesErr := repositories.ReadLibrariesForBook(ctx, bookId)
	if readLibrariesErr != nil {
		return readLibrariesErr
	}

	// Calculate rating avg
	var ratingsSum uint = 0
	librariesCount := len(libraries)
	for _, library := range libraries {
		if library.Rating == 0 {
			librariesCount -= 1
		} else {
			ratingsSum += library.Rating
		}
	}
	var avgRating float32 = 0
	if len(libraries) > 0 {
		avgRating = float32(ratingsSum) / float32(librariesCount)
	}

	// Update rating avg for the given book
	if err := repositories.UpdateBookAvgRating(ctx, bookId, avgRating); err != nil {
		return err
	}

	return nil
}
