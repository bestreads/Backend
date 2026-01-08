package dtos

type ProfileResponse struct {
	UserId               uint64 `json:"userId"`
	Username             string `json:"username"`
	ProfilePicture       string `json:"profilePicture"`
	AccountCreatedAtYear uint   `json:"accountCreatedAtYear"`
	BooksInLibrary       uint   `json:"booksInLibrary"`
	Posts                uint   `json:"posts"`
	Follower             uint   `json:"follower"`
	Following            uint   `json:"following"`
}
