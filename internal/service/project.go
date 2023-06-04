package service

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"time"

	"webportfolio/internal"
)

func (s *service) CreateProject(ctx context.Context, req internal.CreateProjectRequest) (*internal.CreateProjectResponse, error) {
	projectExists, err := s.DB.NewSelect().
		Model((*internal.Project)(nil)).
		Where("name = ?", req.Name).
		Where("user_id = ?", req.UserID).
		Exists(ctx)
	if err != nil {
		return nil, err
	}
	if projectExists {
		return nil, internal.ProjectExistsErr
	}

	project := &internal.Project{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
	}

	_, err = s.DB.NewInsert().Model(project).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &internal.CreateProjectResponse{ID: project.ID}, err
}

func (s *service) GetProject(ctx context.Context, req internal.GetProjectRequest) (*internal.GetProjectResponse, error) {
	project := new(internal.Project)
	query := s.DB.NewSelect().Model(project)
	if req.ID != nil {
		query = query.Where("id = ?", req.ID)
	} else {
		if req.UserID == nil {
			user, err := s.GetUser(ctx, internal.GetUserRequest{
				Name: req.Username,
			})
			if err != nil {
				return nil, err
			}
			req.UserID = &user.ID
		}
		query = query.Where("name = ?", req.Name).Where("user_id = ?", req.UserID)
	}

	err := query.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.ProjectNotFoundErr
		}
		return nil, err
	}
	return &internal.GetProjectResponse{
		ID:            project.ID,
		UserID:        project.UserID,
		Name:          project.Name,
		Description:   project.Description,
		Readme:        project.Readme,
		LikesCount:    project.LikesCount,
		CommentsCount: project.LikesCount,
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}, nil
}

func (s *service) UpdateProject(ctx context.Context, req internal.UpdateProjectRequest) error {
	if req.UserID != 0 {
		project, err := s.GetProject(ctx, internal.GetProjectRequest{
			ID: &req.ID,
		})
		if err != nil {
			return err
		}
		if project.UserID != req.UserID {
			return internal.ProjectNotFoundErr
		}
	}

	query := s.DB.NewUpdate().
		Model((*internal.Project)(nil)).
		Where("id = ?", req.ID)
	if req.Description != nil {
		query = query.Set("description = ?", *req.Description)
	}
	if req.Readme != nil {
		query = query.Set("readme = ?", *req.Readme)
	}
	if req.LikesCount != nil {
		query = query.Set("likes_count = ?", *req.LikesCount)
	}
	if req.UpdatedAt != nil {
		query = query.Set("updated_at = ?", *req.UpdatedAt)
	}
	_, err := query.Exec(ctx)
	return err
}

func (s *service) DeleteProject(ctx context.Context, req internal.DeleteProjectRequest) error {
	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &req.ID,
	})
	if err != nil {
		return err
	}
	if project.UserID != req.UserID {
		return internal.ProjectNotFoundErr
	}

	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID: &req.UserID,
	})
	if err != nil {
		return err
	}

	folder := filepath.Join("content", "projects", user.Username, project.Name)
	os.RemoveAll(folder)

	_, err = s.DB.NewDelete().
		Model((*internal.Project)(nil)).
		Where("id = ?", req.ID).
		Exec(ctx)
	return err
}

func (s *service) UploadProject(ctx context.Context, req internal.UploadProjectRequest) error {
	project, err := s.GetProject(ctx, internal.GetProjectRequest{
		ID: &req.ID,
	})
	if err != nil {
		return err
	}
	if project.UserID != req.UserID {
		return internal.ProjectNotFoundErr
	}
	user, err := s.GetUser(ctx, internal.GetUserRequest{
		ID: &req.UserID,
	})
	if err != nil {
		return err
	}

	archive, err := req.File.Open()
	if err != nil {
		return err
	}
	defer archive.Close()

	folder := filepath.Join("content", "projects", user.Username, project.Name)
	os.RemoveAll(folder)

	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(folder)
		}
	}()

	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, archive)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buff.Bytes())
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		filePath := filepath.Join(folder, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}
	}

	now := time.Now()
	err = s.UpdateProject(ctx, internal.UpdateProjectRequest{
		ID:        req.ID,
		UpdatedAt: &now,
	})
	if err != nil {
		return err
	}

	ok = true
	return nil
}
