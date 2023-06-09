package internal

import (
	"context"
	"mime/multipart"
	"time"
)

type Service interface {
	CreateRelations(context.Context) error

	Login(context.Context, LoginRequest) (*LoginResponse, error)

	CreateUser(context.Context, CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, GetUserRequest) (*GetUserResponse, error)
	UpdateUser(context.Context, UpdateUserRequest) error
	DeleteUser(context.Context, DeleteUserRequest) error
	UploadAvatar(context.Context, UploadAvatarRequest) error

	CreateProject(context.Context, CreateProjectRequest) (*CreateProjectResponse, error)
	GetProject(context.Context, GetProjectRequest) (*GetProjectResponse, error)
	UpdateProject(context.Context, UpdateProjectRequest) error
	DeleteProject(context.Context, DeleteProjectRequest) error
	UploadProject(context.Context, UploadProjectRequest) error

	CreateLike(context.Context, CreateLikeRequest) error
	DeleteLike(context.Context, DeleteLikeRequest) error

	CreateComment(context.Context, CreateCommentRequest) (*CreateCommentResponse, error)
	GetComments(context.Context, GetCommentsRequest) (GetCommentsResponse, error)
	DeleteComment(context.Context, DeleteCommentRequest) error

	// -- unimplemented

	CreateFollow(context.Context, Follow) error
	Following(context.Context, int64) ([]Follow, error)
	Followers(context.Context, int64) ([]Follow, error)
	DeleteFollow(context.Context, int64) error
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`

	Bio string `json:"bio"`
}

type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

type GetUserRequest struct {
	ID       *int64  `query:"id"`
	Name     *string `query:"name"`
	Projects bool    `query:"projects"`

	UserID int64 `json:"-"`
}

type GetUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`

	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Bio            string `json:"bio,omitempty"`
	FollowersCount int64  `json:"followers_count"`

	Projects *[]GetProjectResponse `json:"projects,omitempty"`
}

type UpdateUserRequest struct {
	Fullname *string `json:"fullname"`
	Bio      *string `json:"bio"`

	ID int64 `json:"-"`
}

type DeleteUserRequest struct {
	ID int64 `json:"-"`
}

type UploadAvatarRequest struct {
	File   *multipart.FileHeader
	UserID int64
}

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`

	UserID int64 `json:"-"`
}

type CreateProjectResponse struct {
	ID int64 `json:"id"`
}

type GetProjectRequest struct {
	ID       *int64  `query:"id"`
	Name     *string `query:"name"`
	Username *string `query:"username"`
	UserID   *int64  `query:"user_id"`

	ReqUserID int64 `json:"-"`
}

type GetProjectResponse struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`

	Description string `json:"description"`
	Readme      string `json:"readme"`

	LikesCount    int64 `json:"likes_count"`
	CommentsCount int64 `json:"comments_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsLiked *bool `json:"is_liked,omitempty"`
}

type UpdateProjectRequest struct {
	ID          int64   `json:"id" validate:"required"`
	Description *string `json:"description"`
	Readme      *string `json:"readme"`

	UpdatedAt     *time.Time `json:"-"`
	LikesCount    *int64     `json:"-"`
	CommentsCount *int64     `json:"-"`
	UserID        int64      `json:"-"`
}

type DeleteProjectRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}

type UploadProjectRequest struct {
	ID int64 `form:"id" validate:"required"`

	File   *multipart.FileHeader `json:"-"`
	UserID int64                 `json:"-"`
}

type CreateLikeRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}

type DeleteLikeRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}

type CreateCommentRequest struct {
	ID   int64  `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`

	UserID int64 `json:"-"`
}

type CreateCommentResponse struct {
	ID int64 `json:"id"`
}

type GetCommentsRequest struct {
	ID int64 `query:"id" validate:"required"`
}

type GetCommentsResponse []struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
	User struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Fullname string `json:"fullname"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteCommentRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}
