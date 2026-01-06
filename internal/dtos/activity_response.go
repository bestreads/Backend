package dtos

type ActivityType int

type ActivityResponse[T any] struct {
	Activity T
}
