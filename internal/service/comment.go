package service

import (
	"context"

	"webportfolio/internal"
)

func (*service) CreateComment(context.Context, internal.Comment) error {
	panic("unimplemented")
}

func (*service) Comments(context.Context, int64) ([]internal.Comment, error) {
	panic("unimplemented")
}

func (*service) DeleteComment(context.Context, int64) error {
	panic("unimplemented")
}
