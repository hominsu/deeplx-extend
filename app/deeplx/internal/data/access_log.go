package data

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
	"github.com/oio-network/deeplx-extend/pkgs/pagination"
	"github.com/oio-network/deeplx-extend/schema/ent"
	"github.com/oio-network/deeplx-extend/schema/ent/accesslog"
	"github.com/oio-network/deeplx-extend/schema/ent/user"
)

var _ biz.AccessLogRepo = (*accessLogRepo)(nil)

const accessLogCacheKeyPrefix = "access_log_cache_key_"

type accessLogRepo struct {
	data *Data
	ck   map[string][]string
	sg   *singleflight.Group
	log  *log.Helper
}

// NewAccessLogRepo .
func NewAccessLogRepo(data *Data, logger log.Logger) biz.AccessLogRepo {
	ar := &accessLogRepo{
		data: data,
		sg:   &singleflight.Group{},
		log:  log.NewHelper(log.With(logger, "module", "data/access_log")),
	}
	ar.ck = make(map[string][]string)
	ar.ck["Get"] = []string{"get", "access_log", "id"}
	ar.ck["List"] = []string{"list", "access_log"}
	ar.ck["CountIP"] = []string{"count", "ip"}
	return ar
}

func (r *accessLogRepo) Create(ctx context.Context, log *biz.AccessLog) (*biz.AccessLog, error) {
	m := r.createBuilder(log)
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		al, tErr := toAccessLog(res)
		if tErr != nil {
			return nil, v1.ErrorInternal("internal error: %v", tErr)
		}
		return al, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, v1.ErrorAlreadyExists("access log already exists: %v", err)
	case ent.IsConstraintError(err):
		return nil, v1.ErrorInvalidArgument("invalid argument: %v", err)
	default:
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *accessLogRepo) Get(ctx context.Context, logId int64, view v1.View) (*biz.AccessLog, error) {
	var (
		err error
		key string
		res any
	)
	switch view {
	case v1.View_VIEW_UNSPECIFIED, v1.View_BASIC:
		// key: access_log_cache_key_get_access_log_id:logId
		key = r.cacheKey(strconv.FormatInt(logId, 10), r.ck["Get"]...)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			get := &ent.AccessLog{}
			// get cache
			cErr := r.data.cache.Get(ctx, time.Minute, key, get, func(ctx context.Context) (any, error) { // cache miss
				// get from db
				return r.data.db.AccessLog.Get(ctx, logId)
			})

			return get, cErr
		})
	case v1.View_WITH_EDGE_IDS:
		// key: access_log_cache_key_get_access_log_id_edge_ids:logId
		key = r.cacheKey(strconv.FormatInt(logId, 10), append(r.ck["Get"], "edge_ids")...)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			get := &ent.AccessLog{}
			// get cache
			cErr := r.data.cache.Get(ctx, time.Minute, key, get, func(ctx context.Context) (any, error) { // cache miss
				// get from db
				return r.data.db.AccessLog.Query().
					Where(accesslog.ID(logId)).
					WithOwnerUser(func(query *ent.UserQuery) {
						query.Select(user.FieldID)
						query.Select(user.FieldToken)
					}).
					Only(ctx)
			})

			return get, cErr
		})
	default:
		return nil, v1.ErrorInvalidArgument("invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toAccessLog(res.(*ent.AccessLog))
	case ent.IsNotFound(err): // db miss
		return nil, v1.ErrorNotFound("access log not found: %v", err)
	default: // error
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *accessLogRepo) List(
	ctx context.Context,
	pageSize int,
	pageToken string,
	view v1.View,
) (*biz.AccessLogPage, error) {
	// list access logs
	listQuery := r.data.db.AccessLog.Query().
		Order(ent.Asc(user.FieldID)).
		Limit(pageSize + 1)
	if pageToken != "" {
		token, pErr := pagination.DecodePageToken(pageToken)
		if pErr != nil {
			return nil, v1.ErrorInternal("decode page token err: %v", pErr)
		}
		listQuery = listQuery.Where(accesslog.IDGTE(token))
	}

	// key: access_log_cache_key_list_access_log:pageSize_pageToken
	key := r.cacheKey(
		strings.Join([]string{strconv.FormatInt(int64(pageSize), 10), pageToken}, "_"),
		r.ck["List"]...,
	)

	var (
		err error
		res any
	)

	switch view {
	case v1.View_VIEW_UNSPECIFIED, v1.View_BASIC:
		res, err, _ = r.sg.Do(key, func() (any, error) {
			var ids []int64
			// get cache
			if cErr := r.data.cache.Get(ctx, time.Minute, key, ids, func(ctx context.Context) (any, error) {
				// get from db
				return listQuery.IDs(ctx)
			}); cErr != nil {
				return nil, cErr
			}

			logs := make([]*biz.AccessLog, 0, len(ids))
			for _, id := range ids {
				al, _ := r.Get(ctx, id, view)
				logs = append(logs, al)
			}

			return logs, nil
		})
	case v1.View_WITH_EDGE_IDS:
		res, err, _ = r.sg.Do(key, func() (any, error) {
			var ids []int64
			// get cache
			if cErr := r.data.cache.Get(ctx, time.Minute, key, ids, func(ctx context.Context) (any, error) {
				// get from db
				return listQuery.WithOwnerUser(func(query *ent.UserQuery) {
					query.Select(user.FieldID)
					query.Select(user.FieldToken)
				}).IDs(ctx)
			}); cErr != nil {
				return nil, cErr
			}

			logs := make([]*biz.AccessLog, 0, len(ids))
			for _, id := range ids {
				al, err := r.Get(ctx, id, view)
				if err != nil {
					r.data.log.Warn(err)
				}
				logs = append(logs, al)
			}

			return logs, nil
		})
	default:
		return nil, v1.ErrorInvalidArgument("invalid argument: unknown view")
	}
	switch {
	case err == nil: // db hit, set cache
		logs := res.([]*biz.AccessLog)

		// generate next page token
		var nextPageToken string
		if len(logs) == pageSize+1 {
			nextPageToken, err = pagination.EncodePageToken(logs[len(logs)-1].ID)
			if err != nil {
				return nil, v1.ErrorInternal("encode page token error: %v", err)
			}
			logs = logs[:len(logs)-1]
		}

		return &biz.AccessLogPage{
			AccessLogs:    logs,
			NextPageToken: nextPageToken,
		}, nil
	case ent.IsNotFound(err): // db miss
		return nil, v1.ErrorNotFound("access log not found: %v", err)
	default: // error
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *accessLogRepo) CountIP(ctx context.Context, ip string) (int64, error) {
	// key: access_log_cache_key_count_ip:ip
	key := r.cacheKey(ip, r.ck["CountIP"]...)
	res, err, _ := r.sg.Do(key, func() (any, error) {
		var ids []int64
		// get cache
		cErr := r.data.cache.Get(ctx, time.Minute, key, ids, func(ctx context.Context) (any, error) { // cache miss
			// get from db
			return r.data.db.AccessLog.Query().
				Where(accesslog.IP(ip)).
				IDs(ctx)
		})

		return ids, cErr
	})
	switch {
	case err == nil:
		return int64(len(res.([]int64))), nil
	case ent.IsNotFound(err): // db miss
		return 0, v1.ErrorNotFound("access log not found: %v", err)
	default: // error
		return 0, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *accessLogRepo) createBuilder(log *biz.AccessLog) *ent.AccessLogCreate {
	m := r.data.db.AccessLog.Create()
	m.SetIP(log.IP)
	m.SetCountryName(log.CountryName)
	m.SetCountryCode(log.CountryCode)
	m.SetCreatedAt(log.CreatedAt)
	if log.OwnerUser != nil {
		m.SetOwnerUserID(log.OwnerUser.ID)
	}
	return m
}

func (r *accessLogRepo) cacheKey(unique string, a ...string) string {
	s := strings.Join(a, "_")
	return accessLogCacheKeyPrefix + s + ":" + unique
}

func toAccessLog(e *ent.AccessLog) (*biz.AccessLog, error) {
	a := &biz.AccessLog{}
	a.ID = e.ID
	a.UserID = e.UserID
	a.IP = e.IP
	a.CountryName = e.CountryName
	a.CountryCode = e.CountryCode
	if edg := e.Edges.OwnerUser; edg != nil {
		a.OwnerUser = &biz.User{
			ID: edg.ID,
		}
	}
	return a, nil
}

func toAccessLogList(e []*ent.AccessLog) ([]*biz.AccessLog, error) {
	accessLogList := make([]*biz.AccessLog, len(e))
	for i, entEntity := range e {
		g, err := toAccessLog(entEntity)
		if err != nil {
			return nil, errors.New("convert to accessLogList error")
		}
		accessLogList[i] = g
	}
	return accessLogList, nil
}
