package dtos

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
