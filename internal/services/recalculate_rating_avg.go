package services

import (
	"context"

	"github.com/bestreads/Backend/internal/repositories"
)

// RecalculateRatingAvg gets all ratings of the given book and saves the updated rating avg
func RecalculateRatingAvg(ctx context.Context, bookId uint) error {
	// Get library entries for the given book
	libraries, readLibrariesErr := repositories.ReadLibrariesForBook(ctx, bookId, true)
	if readLibrariesErr != nil {
		return readLibrariesErr
	}

	// Calculate rating avg
	var ratingsSum uint
	var ratedCount uint
	for _, library := range libraries {
		if library.Rating > 0 {
			ratingsSum += library.Rating
			ratedCount++
		}
	}
	var avgRating float32 = 0
	if ratedCount > 0 {
		avgRating = float32(ratingsSum) / float32(ratedCount)
	}

	// Update rating avg for the given book
	if err := repositories.UpdateBookAvgRating(ctx, bookId, avgRating, ratedCount); err != nil {
		return err
	}

	return nil
}
