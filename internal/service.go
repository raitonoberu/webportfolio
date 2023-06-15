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

	CreateAvatar(context.Context, UploadAvatarRequest) error
	// DeleteAvatar(context.Context, DeleteAvatarRequest) error

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

	CreateFollow(context.Context, CreateFollowRequest) error
	GetFollowing(context.Context, GetFollowingRequest) (GetFollowingResponse, error)
	GetFollowers(context.Context, GetFollowersRequest) (GetFollowersResponse, error)
	DeleteFollow(context.Context, DeleteFollowRequest) error

	GetFeed(context.Context, GetFeedRequest) (GetFeedResponse, error)
	GetTrending(context.Context, GetTrendingRequest) (GetTrendingResponse, error)
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

	IsFollowed *bool `json:"is_followed,omitempty"`
}

type UpdateUserRequest struct {
	Fullname *string `json:"fullname"`
	Bio      *string `json:"bio"`

	FollowersCount *int64 `json:"-"`
	ID             int64  `json:"-"`
}

type DeleteUserRequest struct {
	ID int64 `json:"-"`
}

type GetAvatarRequest struct {
	ID int64 `query:"id" validate:"required"`
}

type UploadAvatarRequest struct {
	File   *multipart.FileHeader `json:"-"`
	UserID int64                 `json:"-"`
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
	Folder string `json:"folder"`

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

type CreateFollowRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}

type GetFollowingRequest struct {
	ID int64 `query:"id" validate:"required"`
}

type GetFollowingResponse []struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

type GetFollowersRequest struct {
	ID int64 `query:"id" validate:"required"`
}

type GetFollowersResponse []struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

type DeleteFollowRequest struct {
	ID int64 `json:"id" validate:"required"`

	UserID int64 `json:"-"`
}

type GetFeedRequest struct {
	UserID int64 `json:"-"`
}

type GetFeedResponse []GetProjectResponse

type GetTrendingRequest struct {
	UserID int64 `json:"-"`
}

type GetTrendingResponse []GetProjectResponse
