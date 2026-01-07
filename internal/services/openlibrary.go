package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"resty.dev/v3"
)

const openLibrarySearchURL = "https://openlibrary.org/search.json"
const openLibraryBaseURL = "https://openlibrary.org"

func SearchOpenLibrary(httpClient *resty.Client, ctx context.Context, query string, limit string) error {
	var response dtos.OpenLibraryResponse

	newQuery := strings.ReplaceAll(query, " ", "+")

	params := map[string]string{
		"q":      newQuery,
		"limit":  limit,
		"fields": "isbn,title,author_name,cover_i,first_publish_year,key",
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
			if err := middlewares.DB(ctx).Where("isbn = ?", isbn).First(&existingBook).Error; err == nil {
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

		// Description von Works API holen
		description := fetchWorkDetails(httpClient, ctx, doc.Key)

		book := database.Book{
			Title:       doc.Title,
			Author:      author,
			ISBN:        isbn,
			ReleaseDate: uint64(doc.FirstYear),
			Description: description,
			CoverURL:    coverURL,
		}

		// Neues Buch in DB speichern
		if err := middlewares.DB(ctx).Create(&book).Error; err != nil {
			return fmt.Errorf("failed to save book to database: %w", err)
		}
	}

	return nil
}

func fetchWorkDetails(httpClient *resty.Client, ctx context.Context, workKey string) string {
	if workKey == "" {
		return "Es gibt keine Beschreibung für dieses Buch."
	}

	var workResponse dtos.OpenLibraryWorkResponse

	url := fmt.Sprintf("%s%s.json", openLibraryBaseURL, workKey)

	_, err := httpClient.R().
		SetContext(ctx).
		SetResult(&workResponse).
		Get(url)

	if err != nil {
		return ""
	}

	return extractDescription(workResponse.Description)
}

func extractDescription(desc any) string {
	if desc == nil {
		return ""
	}

	// Wenn es ein String ist
	if str, ok := desc.(string); ok {
		return str
	}

	// Wenn es ein Objekt mit "value" Key
	if obj, ok := desc.(map[string]any); ok {
		if value, exists := obj["value"]; exists {
			if str, ok := value.(string); ok {
				return str
			}
		}
	}

	return ""
}
