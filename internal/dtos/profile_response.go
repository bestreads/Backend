package dtos

type ProfileResponse struct {
	UserID               uint   `json:"userId"`
	Username             string `json:"username"`
	ProfilePicture       string `json:"profilePicture"`
	AccountCreatedAtYear uint   `json:"accountCreatedAtYear"`
	BooksInLibrary       uint   `json:"booksInLibrary"`
	Posts                uint   `json:"posts"`
	FollowersCount       uint   `json:"followersCount"`
	FollowingCount       uint   `json:"followingCount"`
}

type OwnProfileResponse struct {
	ProfileResponse
	Email string `json:"email"`
}
