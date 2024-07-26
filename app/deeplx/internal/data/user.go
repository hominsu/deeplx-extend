package data

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/cache/v9"
	"golang.org/x/sync/singleflight"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent/accesslog"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent/user"
	"github.com/oio-network/deeplx-extend/pkgs/pagination"
)

var _ biz.UserRepo = (*userRepo)(nil)

const userCacheKeyPrefix = "user_cache_key_"

type userRepo struct {
	data *Data
	ck   map[string][]string
	sg   *singleflight.Group
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	ur := &userRepo{
		data: data,
		sg:   &singleflight.Group{},
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
	ur.ck = make(map[string][]string)
	ur.ck["Get"] = []string{"get", "user", "id"}
	return ur
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	m := r.data.db.User.Create()
	m.SetToken(user.Token)
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		u, tErr := toUser(res)
		if tErr != nil {
			return nil, v1.ErrorInternal("internal error: %v", tErr)
		}
		return u, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, v1.ErrorAlreadyExists("user already exists: %v", err)
	case ent.IsConstraintError(err):
		return nil, v1.ErrorInvalidArgument("invalid argument: %v", err)
	default:
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *userRepo) Get(ctx context.Context, userId int64, view v1.View) (*biz.User, error) {
	var (
		err error
		key string
		res any
	)
	switch view {
	case v1.View_VIEW_UNSPECIFIED, v1.View_BASIC:
		// key: user_cache_key_get_user_id:userId
		key = r.cacheKey(strconv.FormatInt(userId, 10), r.ck["Get"]...)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			get := &ent.User{}
			// get cache
			cErr := r.data.cache.Get(ctx, key, get)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				get, cErr = r.data.db.User.Get(ctx, userId)
			}
			return get, cErr
		})
	case v1.View_WITH_EDGE_IDS:
		// key: user_cache_key_get_user_id_edge_ids:userId
		key = r.cacheKey(strconv.FormatInt(userId, 10), append(r.ck["Get"], "edge_ids")...)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			get := &ent.User{}
			// get cache
			cErr := r.data.cache.Get(ctx, key, get)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				get, cErr = r.data.db.User.Query().
					Where(user.ID(userId)).
					WithAccessLogs(func(query *ent.AccessLogQuery) {
						query.Select(accesslog.FieldID)
						query.Select(accesslog.FieldIP)
						query.Select(accesslog.FieldCountryName)
						query.Select(accesslog.FieldCountryCode)
					}).
					Only(ctx)
			}
			return get, cErr
		})
	default:
		return nil, v1.ErrorInvalidArgument("invalid argument: unknown view")
	}
	switch {
	case err == nil: // db hit, set cache
		if err = r.data.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: res.(*ent.User),
			TTL:   r.data.conf.Cache.Ttl.AsDuration(),
		}); err != nil {
			r.log.Errorf("cache error: %v", err)
		}
		return toUser(res.(*ent.User))
	case ent.IsNotFound(err): // db miss
		return nil, v1.ErrorNotFound("user not found: %v", err)
	default: // error
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *userRepo) List(
	ctx context.Context,
	pageSize int,
	pageToken string,
	view v1.View,
) (*biz.UserPage, error) {
	// list users
	listQuery := r.data.db.User.Query().
		Order(ent.Asc(user.FieldID)).
		Limit(pageSize + 1)
	if pageToken != "" {
		token, pErr := pagination.DecodePageToken(pageToken)
		if pErr != nil {
			return nil, v1.ErrorInternal("decode page token err: %v", pErr)
		}
		listQuery = listQuery.Where(user.IDGTE(token))
	}

	var (
		err error
		key string
		res any
	)

	switch view {
	case v1.View_VIEW_UNSPECIFIED, v1.View_BASIC:
		// key: user_cache_key_list_user:pageSize_pageToken
		key = r.cacheKey(
			strings.Join([]string{strconv.FormatInt(int64(pageSize), 10), pageToken}, "_"),
			r.ck["List"]...,
		)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			var entList []*ent.User
			// get cache
			cErr := r.data.cache.GetSkippingLocalCache(ctx, key, &entList)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				entList, cErr = listQuery.All(ctx)
			}
			return entList, cErr
		})
	case v1.View_WITH_EDGE_IDS:
		// key: user_cache_key_list_user_edge_ids:pageSize_pageToken
		key = r.cacheKey(
			strings.Join([]string{strconv.FormatInt(int64(pageSize), 10), pageToken}, "_"),
			append(r.ck["List"], "edge_ids")...,
		)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			var entList []*ent.User
			// get cache
			cErr := r.data.cache.GetSkippingLocalCache(ctx, key, &entList)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				entList, cErr = listQuery.WithAccessLogs(func(query *ent.AccessLogQuery) {
					query.Select(accesslog.FieldID)
					query.Select(accesslog.FieldIP)
					query.Select(accesslog.FieldCountryName)
					query.Select(accesslog.FieldCountryCode)
				}).All(ctx)
			}
			return entList, cErr
		})
	default:
		return nil, v1.ErrorInvalidArgument("invalid argument: unknown view")
	}
	switch {
	case err == nil: // db hit, set cache
		entList := res.([]*ent.User)
		if err = r.data.cache.Set(&cache.Item{
			Ctx:            ctx,
			Key:            key,
			Value:          entList,
			TTL:            r.data.conf.Cache.Ttl.AsDuration(),
			SkipLocalCache: true,
		}); err != nil {
			r.log.Errorf("cache error: %v", err)
		}

		// generate next page token
		var nextPageToken string
		if len(entList) == pageSize+1 {
			nextPageToken, err = pagination.EncodePageToken(entList[len(entList)-1].ID)
			if err != nil {
				return nil, v1.ErrorInternal("encode page token error: %v", err)
			}
			entList = entList[:len(entList)-1]
		}

		userList, tErr := toUserList(entList)
		if tErr != nil {
			return nil, v1.ErrorInternal("internal error: %v", tErr)
		}
		return &biz.UserPage{
			Users:         userList,
			NextPageToken: nextPageToken,
		}, nil
	case ent.IsNotFound(err): // db miss
		return nil, v1.ErrorNotFound("user not found: %v", err)
	default: // error
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *userRepo) cacheKey(unique string, a ...string) string {
	s := strings.Join(a, "_")
	return userCacheKeyPrefix + s + ":" + unique
}

func toUser(e *ent.User) (*biz.User, error) {
	u := &biz.User{}
	u.ID = e.ID
	u.Token = e.Token
	u.CreatedAt = e.CreatedAt
	u.UpdatedAt = e.UpdatedAt
	for _, edg := range e.Edges.AccessLogs {
		u.AccessLogs = append(u.AccessLogs, &biz.AccessLog{
			ID:          edg.ID,
			UserID:      edg.UserID,
			IP:          edg.IP,
			CountryName: edg.CountryName,
			CountryCode: edg.CountryCode,
			CreatedAt:   time.Time{},
		})
	}
	return u, nil
}

func toUserList(e []*ent.User) ([]*biz.User, error) {
	userList := make([]*biz.User, len(e))
	for i, entEntity := range e {
		u, err := toUser(entEntity)
		if err != nil {
			return nil, errors.New("convert to userList error")
		}
		userList[i] = u
	}
	return userList, nil
}
