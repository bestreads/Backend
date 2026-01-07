package dtos

type OpenLibraryBook struct {
	Title      string   `json:"title"`
	AuthorName []string `json:"author_name"`
	ISBN       []string `json:"isbn"`
	FirstYear  int      `json:"first_publish_year"`
	CoverID    int      `json:"cover_i"`
}

type OpenLibraryResponse struct {
	Docs []OpenLibraryBook `json:"docs"`
}
