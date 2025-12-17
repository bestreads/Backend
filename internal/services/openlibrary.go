package services

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"resty.dev/v3"
)

type OpenLibraryBook struct {
	Title      string   `json:"title"`
	AuthorName []string `json:"author_name"`
	ISBN       []string `json:"isbn"`
	FirstYear  int      `json:"first_publish_year"`
}

type OpenLibraryResponse struct {
	Docs []OpenLibraryBook `json:"docs"`
}

func SearchOpenLibrary(httpClient *resty.Client, ctx context.Context, query string) ([]database.Book, error) {
	var response OpenLibraryResponse

	_, err := httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"q":     query,
			"limit": "10",
		}).
		SetResult(&response).
		Get("https://openlibrary.org/search.json")

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
