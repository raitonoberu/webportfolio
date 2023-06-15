package service

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strconv"

	"webportfolio/internal"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(ctx context.Context, req internal.CreateUserRequest) (*internal.CreateUserResponse, error) {
	usernameExists, err := s.DB.NewSelect().
		Model((*internal.User)(nil)).
		Where("username = ?", req.Username).
		Exists(ctx)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, internal.UsernameExistsErr
	}

	emailExists, err := s.DB.NewSelect().
		Model((*internal.User)(nil)).
		Where("email = ?", req.Email).
		Exists(ctx)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, internal.EmailExistsErr
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &internal.User{
		Username: req.Username,
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(passwordHash),
	}

	_, err = s.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}

	token, err := s.newToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &internal.CreateUserResponse{
		ID:    user.ID,
		Token: token,
	}, nil
}

func (s *service) GetUser(ctx context.Context, req internal.GetUserRequest) (*internal.GetUserResponse, error) {
	user := new(internal.User)
	query := s.DB.NewSelect().Model(user)
	if req.ID != nil {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("username = ?", req.Name)
	}
	if req.Projects {
		query = query.Relation("Projects", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.OrderExpr("id DESC")
		})
	}

	err := query.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.UserNotFoundErr
		}
		return nil, err
	}

	result := &internal.GetUserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Fullname:       user.Fullname,
		Email:          user.Email,
		Bio:            user.Bio,
		FollowersCount: user.FollowersCount,
	}
	if req.Projects {
		projects := make([]internal.GetProjectResponse, len(user.Projects))
		for i := 0; i < len(user.Projects); i++ {
			var isLiked *bool
			if req.UserID != 0 {
				l, err := s.isLiked(ctx, user.Projects[i].ID, req.UserID)
				if err != nil {
					return nil, err
				}
				isLiked = &l
			}

			projects[i] = internal.GetProjectResponse{
				ID:            user.Projects[i].ID,
				UserID:        user.Projects[i].UserID,
				Name:          user.Projects[i].Name,
				Description:   user.Projects[i].Description,
				Readme:        user.Projects[i].Readme,
				LikesCount:    user.Projects[i].LikesCount,
				CommentsCount: user.Projects[i].CommentsCount,
				CreatedAt:     user.Projects[i].CreatedAt,
				UpdatedAt:     user.Projects[i].UpdatedAt,

				IsLiked: isLiked,
			}
		}
		result.Projects = &projects
	}
	if req.UserID != 0 {
		followed, err := s.isFollowed(ctx, *req.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		result.IsFollowed = &followed
	}

	return result, nil
}

func (s *service) UpdateUser(ctx context.Context, req internal.UpdateUserRequest) error {
	query := s.DB.NewUpdate().
		Model((*internal.User)(nil)).
		Where("id = ?", req.ID)
	if req.Fullname != nil {
		query = query.Set("fullname = ?", *req.Fullname)
	}
	if req.Bio != nil {
		query = query.Set("bio = ?", *req.Bio)
	}
	if req.FollowersCount != nil {
		query = query.Set("followers_count = ?", *req.FollowersCount)
	}

	_, err := query.Exec(ctx)
	return err
}

func (s *service) DeleteUser(ctx context.Context, req internal.DeleteUserRequest) error {
	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID:       &req.ID,
		Projects: true,
	})
	if err != nil {
		return err
	}

	for _, p := range *user.Projects {
		os.RemoveAll(filepath.Join("content", "projects", p.Folder))
	}
	os.RemoveAll(filepath.Join("content", "avatars", strconv.FormatInt(user.ID, 10)))

	// TODO: delete likes & comments & follows

	_, err = s.DB.NewDelete().
		Model((*internal.Project)(nil)).
		Where("user_id = ?", req.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = s.DB.NewDelete().
		Model((*internal.User)(nil)).
		Where("id = ?", req.ID).
		Exec(ctx)
	return err
}
