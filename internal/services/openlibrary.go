package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"resty.dev/v3"
)

const openLibrarySearchURL = "https://openlibrary.org/search.json" // Open Library Search-Endpoint
const openLibraryBaseURL = "https://openlibrary.org"               // Basis-URL für Work-Details und Beschreibungen
const maxConcurrentRequests = 7                                    // Rate limiting: maximal 7 gleichzeitige Requests

// SearchOpenLibrary führt eine OpenLibrary-Suche aus, lädt parallel zugehörige Work-Descriptions
// und speichert die gefundenen Bücher inklusive Beschreibung und Cover-URL in der Datenbank.
func SearchOpenLibrary(httpClient *resty.Client, ctx context.Context, query string, limit string) error {
	var response dtos.OpenLibraryResponse

	_, err := httpClient.R().
		SetContext(ctx).
		SetQueryParams(buildQuery(query, limit)).
		SetResult(&response).
		Get(openLibrarySearchURL)
	if err != nil {
		return err
	}

	var (
		res   []database.Book
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	for i, book := range response.Docs {
		wg.Add(1)
		go func(i int, book dtos.OpenLibraryBook) {
			defer wg.Done()

			// hier holen wir die daten, seperat vom sperren
			single, err := metadataSingle(httpClient, ctx, book)
			if err != nil {
				log := middlewares.Logger(ctx)
				log.Err(err).Msg(fmt.Sprintf("worker %d returned an error", i))
				// einfach "leise" abbrechen
				return
			}

			// lsp kaputt?
			lsp_workaround := dtos.Olibrary2book(single)

			// hier ist die array-sperrlogik
			mutex.Lock()
			defer mutex.Unlock()
			// das ist kein error
			res = append(res, lsp_workaround)

		}(i, book)
	}

	wg.Wait()

	if err := insertNewBooks(ctx, res); err != nil {
		return err
	}

	return nil
}

func insertNewBooks(ctx context.Context, books []database.Book) error {
	// das hier ist immer noch xddd
	err := middlewares.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, book := range books {
			if book.ISBN != "" {
				if err := tx.
					Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "isbn"}}, DoNothing: true}).
					Create(&book).Error; err != nil {
					return err
				}

			} else {
				if err := tx.Create(&book).Error; err != nil {
					return err
				}

			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func metadataSingle(client *resty.Client, ctx context.Context, book dtos.OpenLibraryBook) (dtos.OlibFullData, error) {
	isbn, err := dtos.UnwrapFirst(book.ISBN)
	if err != nil {
		return dtos.OlibFullData{}, err
	}

	author, err := dtos.UnwrapFirst(book.AuthorName)
	if err != nil {
		return dtos.OlibFullData{}, err
	}

	desc := fetchWorkDetails(client, ctx, book.Key)

	cacheId, err := database.CacheMedia(coverUrlfmt(book.CoverID))
	if err != nil {
		return dtos.OlibFullData{}, err
	}
	cachedURL := CacheKey2Url(ctx, cacheId)

	return dtos.OlibFullData{
		Title:       book.Title,
		Author:      author,
		ISBN:        isbn,
		Year:        book.FirstYear,
		CoverURL:    cachedURL,
		Description: desc,
	}, nil
}

func coverUrlfmt(id int) string {
	return fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", id)

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

func buildQuery(query string, limit string) map[string]string {
	urlfmt := strings.ReplaceAll(query, " ", "+")
	return map[string]string{
		"q":      "title_suggest:\"" + urlfmt + "\"",
		"limit":  limit,
		"fields": "isbn,title,author_name,cover_i,first_publish_year,key",
	}

}

func CacheKey2Url(ctx context.Context, id uint64) string {
	cfg := middlewares.Config(ctx)
	return fmt.Sprintf("%s://%s%s/v1/media/%d", cfg.ApiProtocol, cfg.ApiDomain, cfg.ApiBasePath, id)
}
