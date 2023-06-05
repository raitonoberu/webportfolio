package internal

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID       int64  `bun:"id,pk,autoincrement"`
	Email    string `bun:"email,notnull,unique"`
	Username string `bun:"username,notnull,unique"`
	Password string `bun:"password,notnull"`

	Fullname string `bun:"fullname,notnull"`
	Bio      string `bun:"bio"`
	Avatar   bool   `bun:"avatar"`

	FollowersCount int64    `bun:"followers_count"`
	Followers      []Follow `bun:"rel:has-many,join:id=following_id"`
	Following      []Follow `bun:"rel:has-many,join:id=follower_id"`

	Projects []Project `bun:"rel:has-many,join:id=user_id"`
}

type Project struct {
	bun.BaseModel

	ID     int64 `bun:"id,pk,autoincrement"`
	UserID int64 `bun:"user_id,notnull"`

	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	Likes         []Like    `bun:"rel:has-many,join:id=user_id"`
	LikesCount    int64     `bun:"likes_count"`
	Comments      []Comment `bun:"rel:has-many,join:id=user_id"`
	CommentsCount int64     `bun:"comments_count"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`

	Readme string `bun:"readme"`
}

type Like struct {
	bun.BaseModel

	ID        int64 `bun:"id,pk,autoincrement"`
	UserID    int64 `bun:"user_id,notnull"`
	ProjectID int64 `bun:"project_id,notnull"`
}

type Comment struct {
	bun.BaseModel

	ID        int64     `bun:"id,pk,autoincrement"`
	UserID    int64     `bun:"user_id,notnull"`
	ProjectID int64     `bun:"project_id,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`

	Text string `bun:"text"`

	User *User `bun:"rel:has-one,join:user_id=id"`
}

type Follow struct {
	bun.BaseModel

	ID          int64 `bun:"id,pk,autoincrement"`
	FollowerID  int64 `bun:"follower_id,notnull"`
	FollowingID int64 `bun:"following_id,notnull"`
}
