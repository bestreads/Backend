package services

import (
	"context"

	"github.com/bestreads/Backend/internal/dtos"
	"github.com/bestreads/Backend/internal/repositories"
)

const CONTENT_LIMIT uint = 10

func GetActivity[T any](ctx context.Context, uids []uint) ([]dtos.ActivityResponse[T], error) {
	posts, err := repositories.GetPostActivity(ctx, uids, CONTENT_LIMIT)

	ret := make([]dtos.ActivityResponse[T], CONTENT_LIMIT*2)

	if err != nil {
		return []dtos.ActivityResponse[T]{}, err
	}

	converted, err := ConvertPost(posts)
	if err != nil {
		return []dtos.ActivityResponse[T]{}, err
	}

	first := wrap(converted)

	activity, err := repositories.GetBookActivity(ctx, uids, CONTENT_LIMIT)
	if err != nil {
		return []dtos.ActivityResponse[T]{}, err
	}
	first = append(first, wrap(activity))

}

func wrap[T any](data []T) []dtos.ActivityResponse[T] {
	ret := make([]dtos.ActivityResponse[T], len(data))
	for idx, elem := range data {
		ret[idx] = dtos.ActivityResponse[T]{Activity: elem}
	}
	return ret
}
