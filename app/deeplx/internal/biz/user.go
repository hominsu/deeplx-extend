package biz

import (
	"context"
	"time"

	"github.com/google/uuid"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
)

type User struct {
	ID         int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Token      uuid.UUID
	AccessLogs []*AccessLog
}

type UserPage struct {
	Users         []*User
	NextPageToken string
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, userId int64, view v1.View) (*User, error)
	List(ctx context.Context, pageSize int, pageToken string, view v1.View) (*UserPage, error)
}
