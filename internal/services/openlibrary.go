package services

import (
	"context"
	"fmt"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"gorm.io/gorm"
	"resty.dev/v3"
)

const openLibrarySearchURL = "https://openlibrary.org/search.json"

func SearchOpenLibrary(httpClient *resty.Client, db *gorm.DB, ctx context.Context, query string, limit string) error {
	var response dtos.OpenLibraryResponse

	params := map[string]string{
		"q":      query,
		"limit":  limit,
		"fields": "isbn,title,author_name,cover_i,first_publish_year,subject",
	}

	_, err := httpClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&response).
		Get(openLibrarySearchURL)

	if err != nil {
		return err
	}

	for _, doc := range response.Docs {
		isbn := ""
		if len(doc.ISBN) > 0 {
			isbn = doc.ISBN[0]
		}

		// Prüfen ob Buch mit dieser ISBN bereits existiert
		if isbn != "" {
			var existingBook database.Book
			if err := db.WithContext(ctx).Where("isbn = ?", isbn).First(&existingBook).Error; err == nil {
				// Buch existiert bereits, überspringen
				continue
			}
		}

		author := ""
		if len(doc.AuthorName) > 0 {
			author = doc.AuthorName[0]
		}

		coverURL := ""
		if doc.CoverID > 0 {
			coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", doc.CoverID)
		}

		book := database.Book{
			Title:       doc.Title,
			Author:      author,
			ISBN:        isbn,
			ReleaseDate: uint64(doc.FirstYear),
			CoverURL:    coverURL,
		}

		// Neues Buch in DB speichern
		if err := db.WithContext(ctx).Create(&book).Error; err != nil {
			return fmt.Errorf("failed to save book to database: %w", err)
		}
	}

	return nil
}
