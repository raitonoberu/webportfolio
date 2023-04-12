package service

import (
	"context"

	"webportfolio/internal"
)

func (*service) CreateFollow(context.Context, internal.Follow) error {
	panic("unimplemented")
}

func (*service) Followers(context.Context, int64) ([]internal.Follow, error) {
	panic("unimplemented")
}

func (*service) Following(context.Context, int64) ([]internal.Follow, error) {
	panic("unimplemented")
}

func (*service) DeleteFollow(context.Context, int64) error {
	panic("unimplemented")
}
