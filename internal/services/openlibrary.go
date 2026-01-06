package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/dtos"
	"resty.dev/v3"
)

func SearchOpenLibrary(httpClient *resty.Client, ctx context.Context, query string, limit string) ([]database.Book, error) {
	var response dtos.OpenLibraryResponse
	var err error
	if limit != "" {
		_, err = httpClient.R().
			SetContext(ctx).
			SetQueryParams(map[string]string{
				"q":     query,
				"limit": limit,
			}).
			SetResult(&response).
			Get("https://openlibrary.org/search.json")
	} else {
		_, err = httpClient.R().
			SetContext(ctx).
			SetQueryParams(map[string]string{
				"q": query,
			}).
			SetResult(&response).
			Get("https://openlibrary.org/search.json")
	}

	if err != nil {
		return nil, err
	}

	var books []database.Book
	for _, doc := range response.Docs {
		isbn := ""
		if len(doc.ISBN) > 0 {
			isbn = doc.ISBN[0]
		}

		author := ""
		if len(doc.AuthorName) > 0 {
			author = doc.AuthorName[0]
		}

		book := database.Book{
			Title:       doc.Title,
			Author:      author,
			ISBN:        isbn,
			ReleaseDate: uint64(doc.FirstYear),
		}
		books = append(books, book)
	}

	return books, nil
}
