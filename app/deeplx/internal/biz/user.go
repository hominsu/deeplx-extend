package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
)

type User struct {
	ID         int64
	Token      uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
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

func ToUser(p *v1.User) (*User, error) {
	token, err := uuid.Parse(p.GetToken())
	if err != nil {
		return nil, err
	}

	u := &User{}
	u.ID = p.GetId()
	u.Token = token
	u.CreatedAt = p.GetCreatedAt().AsTime()
	u.UpdatedAt = p.GetUpdatedAt().AsTime()
	for _, log := range p.GetAccessLogs() {
		accessLog, err := ToAccessLog(log)
		if err != nil {
			return nil, err
		}
		u.AccessLogs = append(u.AccessLogs, accessLog)
	}

	return u, nil
}

func ToUserList(p []*v1.User) ([]*User, error) {
	userList := make([]*User, 0, len(p))
	for _, pbEntity := range p {
		u, err := ToUser(pbEntity)
		if err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}
	return userList, nil
}

func ToProtoUser(u *User) (*v1.User, error) {
	p := &v1.User{
		Id:         u.ID,
		Token:      u.Token.String(),
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
		AccessLogs: make([]*v1.AccessLog, 0, len(u.AccessLogs)),
	}
	for _, log := range u.AccessLogs {
		pl, err := ToProtoAccessLog(log)
		if err != nil {
			return nil, err
		}
		p.AccessLogs = append(p.AccessLogs, pl)
	}
	return p, nil
}

func ToProtoUserList(u []*User) ([]*v1.User, error) {
	pbList := make([]*v1.User, 0, len(u))
	for _, userEntity := range u {
		pbUser, err := ToProtoUser(userEntity)
		if err != nil {
			return nil, err
		}
		pbList = append(pbList, pbUser)
	}

	return pbList, nil
}
