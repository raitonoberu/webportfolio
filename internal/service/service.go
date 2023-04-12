package service

import (
	"context"

	"webportfolio/internal"

	"github.com/uptrace/bun"
)

type service struct {
	DB     *bun.DB
	Secret []byte
}

func New(db *bun.DB, secret string) internal.Service {
	return &service{db, []byte(secret)}
}

func (s *service) CreateRelations(ctx context.Context) error {
	models := []interface{}{
		(*internal.User)(nil),
		(*internal.Project)(nil),
		(*internal.Like)(nil),
		(*internal.Comment)(nil),
		(*internal.Follow)(nil),
	}

	for _, m := range models {
		_, err := s.DB.NewCreateTable().
			IfNotExists().
			Model(m).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
