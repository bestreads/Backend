package dtos

import (
	"errors"

	"github.com/bestreads/Backend/internal/database"
)

type OpenLibraryBook struct {
	Title      string   `json:"title"`
	AuthorName []string `json:"author_name"`
	ISBN       []string `json:"isbn"`
	FirstYear  int      `json:"first_publish_year"`
	CoverID    int      `json:"cover_i"`
	Key        string   `json:"key"`
}

type OpenLibraryResponse struct {
	Docs []OpenLibraryBook `json:"docs"`
}

type OpenLibraryWorkResponse struct {
	Description any `json:"description"`
}

type OlibFullData struct {
	Title       string
	Author      string
	ISBN        string
	Year        int
	CoverURL    string
	Description string
}

func Olibrary2book(olibbook OlibFullData) database.Book {

	book := database.Book{
		Title:       olibbook.Title,
		Author:      olibbook.Author,
		ISBN:        olibbook.ISBN,
		ReleaseDate: uint64(olibbook.Year),
		Description: olibbook.Description,
		CoverURL:    olibbook.CoverURL,
	}

	return book
}

func UnwrapFirst[T any](array []T) (T, error) {
	if len(array) > 0 {
		return array[0], nil
	}
	return *new(T), errors.New("No element found")
}
