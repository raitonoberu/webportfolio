package service

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"webportfolio/internal"

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
		query = query.Relation("Projects")
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
		Avatar:         user.Avatar,
		FollowersCount: user.FollowersCount,
	}
	if req.Projects {
		projects := make([]internal.GetProjectResponse, len(user.Projects))
		for i := 0; i < len(user.Projects); i++ {
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
			}
		}
		result.Projects = &projects
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
	if req.Avatar != nil {
		query = query.Set("avatar = ?", *req.Avatar)
	}

	_, err := query.Exec(ctx)
	return err
}

func (s *service) DeleteUser(ctx context.Context, req internal.DeleteUserRequest) error {
	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID: &req.ID,
	})
	if err != nil {
		return err
	}

	os.RemoveAll(filepath.Join("content", "projects", user.Username))
	os.RemoveAll(filepath.Join("content", "avatars", strconv.FormatInt(user.ID, 10)))

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

func (s *service) UploadAvatar(ctx context.Context, req internal.UploadAvatarRequest) error {
	img, err := req.File.Open()
	if err != nil {
		return err
	}
	defer img.Close()

	// TODO: resize/convert?
	folder := filepath.Join("content", "avatars")
	os.MkdirAll(folder, os.ModePerm)
	userID := strconv.FormatInt(req.UserID, 10)

	dst, err := os.Create(filepath.Join(folder, userID))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, img); err != nil {
		return err
	}

	avatar := true
	return s.UpdateUser(ctx, internal.UpdateUserRequest{
		ID: req.UserID, Avatar: &avatar,
	})
}
