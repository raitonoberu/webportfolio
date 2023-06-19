package service

import (
	"context"
	"webportfolio/internal"

	"github.com/uptrace/bun"
)

func (s *service) GetFeed(ctx context.Context, req internal.GetFeedRequest) (internal.GetFeedResponse, error) {
	var user_ids []int64

	err := s.DB.NewSelect().
		Model((*internal.Follow)(nil)).
		Where("follower_id = ?", req.UserID).
		Column("following_id").
		Scan(ctx, &user_ids)
	if err != nil {
		return nil, err
	}
	if len(user_ids) == 0 {
		return nil, nil
	}

	var projects []internal.Project

	err = s.DB.NewSelect().
		Model(&projects).
		Where("user_id IN (?)", bun.In(user_ids)).
		OrderExpr("updated_at DESC").
		Limit(20).
		Scan(ctx)

	result := make(internal.GetFeedResponse, len(projects))
	for i := 0; i < len(projects); i++ {
		var isLiked *bool
		if req.UserID != 0 {
			l, err := s.isLiked(ctx, projects[i].ID, req.UserID)
			if err != nil {
				return nil, err
			}
			isLiked = &l
		}

		result[i] = internal.GetProjectResponse{
			ID:            projects[i].ID,
			UserID:        projects[i].UserID,
			Name:          projects[i].Name,
			Folder:        projects[i].Folder,
			Description:   projects[i].Description,
			Readme:        projects[i].Readme,
			LikesCount:    projects[i].LikesCount,
			CommentsCount: projects[i].CommentsCount,
			CreatedAt:     projects[i].CreatedAt,
			UpdatedAt:     projects[i].UpdatedAt,

			IsLiked: isLiked,
		}
	}
	return result, err
}

func (s *service) GetTrending(ctx context.Context, req internal.GetTrendingRequest) (internal.GetTrendingResponse, error) {
	var projects []internal.Project
	err := s.DB.NewSelect().
		Model(&projects).
		OrderExpr("likes_count DESC").
		Limit(20).
		Scan(ctx)

	result := make(internal.GetTrendingResponse, len(projects))
	for i := 0; i < len(projects); i++ {
		// TODO: remove repetition
		var isLiked *bool
		if req.UserID != 0 {
			l, err := s.isLiked(ctx, projects[i].ID, req.UserID)
			if err != nil {
				return nil, err
			}
			isLiked = &l
		}

		result[i] = internal.GetProjectResponse{
			ID:            projects[i].ID,
			UserID:        projects[i].UserID,
			Name:          projects[i].Name,
			Folder:        projects[i].Folder,
			Description:   projects[i].Description,
			Readme:        projects[i].Readme,
			LikesCount:    projects[i].LikesCount,
			CommentsCount: projects[i].CommentsCount,
			CreatedAt:     projects[i].CreatedAt,
			UpdatedAt:     projects[i].UpdatedAt,

			IsLiked: isLiked,
		}
	}
	return result, err
}
