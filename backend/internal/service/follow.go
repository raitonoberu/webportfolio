package service

import (
	"context"

	"webportfolio/internal"
)

func (s *service) CreateFollow(ctx context.Context, req internal.CreateFollowRequest) error {
	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID:     &req.ID,
		UserID: req.UserID,
	})
	if err != nil {
		return err
	}

	if *user.IsFollowed {
		return internal.UserFollowedErr
	}

	follow := &internal.Follow{
		FollowerID:  req.UserID,
		FollowingID: req.ID,
	}
	_, err = s.DB.NewInsert().
		Model(follow).
		Exec(ctx)
	if err != nil {
		return err
	}

	followersCount := user.FollowersCount + 1
	return s.UpdateUser(ctx, internal.UpdateUserRequest{
		ID:             req.ID,
		FollowersCount: &followersCount,
	})
}

func (s *service) GetFollowing(ctx context.Context, req internal.GetFollowingRequest) (internal.GetFollowingResponse, error) {
	var follows []internal.Follow
	err := s.DB.NewSelect().
		Model(&follows).
		Relation("Following").
		Where("follower_id = ?", req.ID).
		OrderExpr("id DESC").
		Scan(ctx)

	result := make(internal.GetFollowingResponse, len(follows))
	for i := 0; i < len(follows); i++ {
		result[i].ID = follows[i].Following.ID
		result[i].Username = follows[i].Following.Username
		result[i].Fullname = follows[i].Following.Fullname
	}
	return result, err
}

func (s *service) GetFollowers(ctx context.Context, req internal.GetFollowersRequest) (internal.GetFollowersResponse, error) {
	var follows []internal.Follow
	err := s.DB.NewSelect().
		Model(&follows).
		Relation("Follower").
		Where("following_id = ?", req.ID).
		OrderExpr("id DESC").
		Scan(ctx)

	result := make(internal.GetFollowersResponse, len(follows))
	for i := 0; i < len(follows); i++ {
		result[i].ID = follows[i].Follower.ID
		result[i].Username = follows[i].Follower.Username
		result[i].Fullname = follows[i].Follower.Fullname
	}
	return result, err
}

func (s *service) DeleteFollow(ctx context.Context, req internal.DeleteFollowRequest) error {
	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID:     &req.ID,
		UserID: req.UserID,
	})
	if err != nil {
		return err
	}

	if !*user.IsFollowed {
		return internal.UserNotFollowedErr
	}

	_, err = s.DB.NewDelete().
		Model((*internal.Follow)(nil)).
		Where("follower_id = ?", req.UserID).
		Where("following_id = ?", req.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	followersCount := user.FollowersCount - 1
	return s.UpdateUser(ctx, internal.UpdateUserRequest{
		ID:             req.ID,
		FollowersCount: &followersCount,
	})
}

func (s *service) isFollowed(ctx context.Context, id, userID int64) (bool, error) {
	return s.DB.NewSelect().
		Model((*internal.Follow)(nil)).
		Where("following_id = ?", id).
		Where("follower_id = ?", userID).
		Exists(ctx)
}
