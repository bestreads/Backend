package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"github.com/bestreads/Backend/internal/repositories"
	"gorm.io/gorm"
	"resty.dev/v3"
)

const openLibrarySearchURL = "https://openlibrary.org/search.json" // Open Library Search-Endpoint
const openLibraryBaseURL = "https://openlibrary.org"               // Basis-URL für Work-Details und Beschreibungen
const maxConcurrentRequests = 15                                   // Rate limiting: maximal 15 gleichzeitige Requests

// SearchOpenLibrary führt eine OpenLibrary-Suche aus, lädt parallel zugehörige Work-Descriptions
// und speichert die gefundenen Bücher inklusive Beschreibung und Cover-URL in der Datenbank.
func SearchOpenLibrary(httpClient *resty.Client, ctx context.Context, query string, limit string) error {
	var response dtos.OpenLibraryResponse

	newQuery := strings.ReplaceAll(query, " ", "+")

	params := map[string]string{
		"q":      "title_suggest:\"" + newQuery + "\" author:\"" + newQuery + "\"",
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

	// Descriptions und Cover-URLs parallel fetchen mit Rate Limiting
	descriptions := fetchAllDescriptions(httpClient, ctx, response.Docs)
	coverURLs := fetchAllCoverURLs(ctx, response.Docs)

	// Transaction für alle Buch-Inserts
	if err := middlewares.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for i, doc := range response.Docs {
			if err := insertBook(tx, ctx, doc, descriptions[i], coverURLs[i]); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// fetchAllCoverURLs holt alle Cover-Bilder parallel und cached sie
func fetchAllCoverURLs(ctx context.Context, docs []dtos.OpenLibraryBook) []string {
	coverURLs := make([]string, len(docs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentRequests) // Rate limiting

	for i, doc := range docs {
		wg.Add(1)
		go func(index int, coverID int) {
			defer wg.Done()

			if coverID <= 0 {
				coverURLs[index] = ""
				return
			}

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			hash, err := database.CacheMedia(fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", coverID))
			if err != nil {
				coverURLs[index] = ""
				return
			}
			coverURLs[index] = fmt.Sprintf("%s/api/v1/media/%d", middlewares.Config(ctx).ApiBaseURL, hash)
		}(i, doc.CoverID)
	}

	wg.Wait()
	return coverURLs
}

// insertBook fügt ein Buch in die Datenbank ein
func insertBook(tx *gorm.DB, ctx context.Context, doc dtos.OpenLibraryBook, description string, cachedURL string) error {
	isbn := ""
	if len(doc.ISBN) > 0 {
		isbn = doc.ISBN[0]
	}

	author := ""
	if len(doc.AuthorName) > 0 {
		author = doc.AuthorName[0]
	}

	if description == "" {
		description = "Es gibt keine Beschreibung für dieses Buch."
	}

	book := database.Book{
		Title:       doc.Title,
		Author:      author,
		ISBN:        isbn,
		ReleaseDate: uint64(doc.FirstYear),
		Description: description,
		CoverURL:    cachedURL,
	}

	// Für Bücher mit ISBN: ON CONFLICT DO NOTHING für idempotentes Verhalten
	// Für Bücher ohne ISBN: normales Create (jedes Buch wird eingefügt)
	if isbn != "" {
		if err := repositories.CreateBookNoISBN(tx, ctx, &book); err != nil {
			return fmt.Errorf("failed to save book to database: %w", err)
		}
	} else {
		if err := repositories.CreateBookISBN(tx, ctx, &book); err != nil {
			return fmt.Errorf("failed to save book to database: %w", err)
		}
	}

	return nil
}

// fetchAllDescriptions holt alle Descriptions parallel mit Rate Limiting
func fetchAllDescriptions(httpClient *resty.Client, ctx context.Context, docs []dtos.OpenLibraryBook) []string {
	descriptions := make([]string, len(docs))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentRequests) // Rate limiting

	for i, doc := range docs {
		wg.Add(1)
		go func(index int, workKey string) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			descriptions[index] = fetchWorkDetails(httpClient, ctx, workKey)
		}(i, doc.Key)
	}

	wg.Wait()
	return descriptions
}

// fetchWorkDetails ruft die Work-JSON von Open Library ab und extrahiert die Description
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

// extractDescription wandelt Description-Objekte oder Strings in einen lesbaren Text um
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
