package dtos

type UpdateReviewRequest struct {
	BookID uint `json:"bookId"`
	Rating uint `json:"rating" validate:"required,min=1,max=5"`
}
