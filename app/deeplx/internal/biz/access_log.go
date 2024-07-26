package biz

import (
	"context"
	"time"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
)

type AccessLog struct {
	ID          int64
	UserID      int64
	IP          string
	CountryName string
	CountryCode string
	CreatedAt   time.Time
	OwnerUser   *User
}

type AccessLogPage struct {
	AccessLogs    []*AccessLog
	NextPageToken string
}

type AccessLogRepo interface {
	Create(ctx context.Context, log *AccessLog) (*AccessLog, error)
	Get(ctx context.Context, logId int64, view v1.View) (*AccessLog, error)
	List(ctx context.Context, pageSize int, pageToken string, view v1.View) (*AccessLogPage, error)
	CountIP(ctx context.Context, ip string) (int64, error)
}
