package service

import (
	"context"
	"database/sql"

	"webportfolio/internal"
)

func (s *service) CreateComment(ctx context.Context, req internal.CreateCommentRequest) (*internal.CreateCommentResponse, error) {
	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &req.ID,
	})
	if err != nil {
		return nil, err
	}

	comment := &internal.Comment{
		UserID:    req.UserID,
		ProjectID: req.ID,
		Text:      req.Text,
	}
	_, err = s.DB.NewInsert().
		Model(comment).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	commentsCount := project.CommentsCount + 1
	return &internal.CreateCommentResponse{
			ID: comment.ID,
		}, s.UpdateProject(ctx, internal.UpdateProjectRequest{
			ID:            req.ID,
			CommentsCount: &commentsCount,
		})
}

func (s *service) GetComments(ctx context.Context, req internal.GetCommentsRequest) (internal.GetCommentsResponse, error) {
	var comments []internal.Comment
	err := s.DB.NewSelect().
		Model(&comments).
		Where("project_id = ?", req.ID).
		Relation("User").
		Scan(ctx)

	result := make(internal.GetCommentsResponse, len(comments))
	for i := 0; i < len(comments); i++ {
		result[i].ID = comments[i].ID
		result[i].Text = comments[i].Text
		result[i].User.ID = comments[i].User.ID
		result[i].User.Username = comments[i].User.Username
		result[i].User.Fullname = comments[i].User.Fullname
		result[i].CreatedAt = comments[i].CreatedAt
	}
	return result, err
}

func (s *service) DeleteComment(ctx context.Context, req internal.DeleteCommentRequest) error {
	comment := new(internal.Comment)
	err := s.DB.NewSelect().
		Model(comment).
		Where("id = ?", req.ID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return internal.CommentNotFoundErr
		}
		return err
	}

	if comment.UserID != req.UserID {
		return internal.CommentNotFoundErr
	}

	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &comment.ProjectID,
	})
	if err != nil {
		return err
	}

	_, err = s.DB.NewDelete().
		Model(comment).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	commentsCount := project.CommentsCount - 1
	return s.UpdateProject(ctx, internal.UpdateProjectRequest{
		ID:            comment.ProjectID,
		CommentsCount: &commentsCount,
	})
}
