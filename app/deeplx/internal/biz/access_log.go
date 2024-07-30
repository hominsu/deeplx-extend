package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"

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

type AccessLogUseCase struct {
	repo AccessLogRepo
	log  *log.Helper
}

func NewAccessLogUsecase(repo AccessLogRepo, logger log.Logger) *AccessLogUseCase {
	return &AccessLogUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz")),
	}
}

func (uc *AccessLogUseCase) Create(ctx context.Context, log *v1.AccessLog) (*v1.AccessLog, error) {
	l, err := ToAccessLog(log)
	if err != nil {
		return nil, err
	}

	ret, err := uc.repo.Create(ctx, l)
	if err != nil {
		return nil, err
	}

	return ToProtoAccessLog(ret)
}

func ToAccessLog(p *v1.AccessLog) (*AccessLog, error) {
	al := &AccessLog{
		ID:          p.GetId(),
		UserID:      p.GetUserId(),
		IP:          p.GetIp(),
		CountryName: p.GetCountryName(),
		CountryCode: p.GetCountryCode(),
		CreatedAt:   p.GetCreatedAt().AsTime(),
	}
	if p.OwnerUser != nil {
		al.OwnerUser = &User{
			ID: p.OwnerUser.GetId(),
		}
	}

	return al, nil
}

func ToAccessLogList(p []*v1.AccessLog) ([]*AccessLog, error) {
	alList := make([]*AccessLog, 0, len(p))
	for _, pbEntity := range p {
		al, err := ToAccessLog(pbEntity)
		if err != nil {
			return nil, err
		}
		alList = append(alList, al)
	}
	return alList, nil
}

func ToProtoAccessLog(al *AccessLog) (*v1.AccessLog, error) {
	p := &v1.AccessLog{
		Id:          al.ID,
		UserId:      al.UserID,
		Ip:          al.IP,
		CountryName: al.CountryName,
		CountryCode: al.CountryCode,
		CreatedAt:   timestamppb.New(al.CreatedAt),
	}
	if al.OwnerUser != nil {
		p.OwnerUser = &v1.User{
			Id: al.OwnerUser.ID,
		}
	}

	return p, nil
}

func ToProtoAccessLogList(alList []*AccessLog) ([]*v1.AccessLog, error) {
	pbList := make([]*v1.AccessLog, 0, len(alList))
	for _, accessLogEntity := range alList {
		pbAccessLog, err := ToProtoAccessLog(accessLogEntity)
		if err != nil {
			return nil, err
		}
		pbList = append(pbList, pbAccessLog)
	}
	return pbList, nil
}
