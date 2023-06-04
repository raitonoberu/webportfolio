package service

import (
	"context"
	"database/sql"
	"webportfolio/internal"
)

func (s *service) CreateLike(ctx context.Context, req internal.CreateLikeRequest) error {
	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &req.ID,
	})
	if err != nil {
		return err
	}

	likeExists, err := s.DB.NewSelect().
		Model((*internal.Like)(nil)).
		Where("project_id = ?", req.ID).
		Where("user_id = ?", req.UserID).
		Exists(ctx)
	if err != nil {
		return err
	}
	if likeExists {
		return internal.ProjectLikedErr
	}

	like := &internal.Like{
		ProjectID: req.ID,
		UserID:    req.UserID,
	}
	_, err = s.DB.NewInsert().
		Model(like).
		Exec(ctx)
	if err != nil {
		return err
	}

	likesCount := project.LikesCount + 1
	return s.UpdateProject(ctx, internal.UpdateProjectRequest{
		ID:         req.ID,
		LikesCount: &likesCount,
	})
}

func (s *service) DeleteLike(ctx context.Context, req internal.DeleteLikeRequest) error {
	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &req.ID,
	})
	if err != nil {
		return err
	}

	like := new(internal.Like)
	err = s.DB.NewSelect().
		Model(like).
		Where("project_id = ?", req.ID).
		Where("user_id = ?", req.UserID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return internal.ProjectNotLikedErr
		}
		return err
	}

	_, err = s.DB.NewDelete().
		Model(like).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	likesCount := project.LikesCount - 1
	return s.UpdateProject(ctx, internal.UpdateProjectRequest{
		ID:         req.ID,
		LikesCount: &likesCount,
	})
}
