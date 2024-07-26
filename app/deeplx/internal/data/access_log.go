package data

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/cache/v9"
	"golang.org/x/sync/singleflight"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent/accesslog"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent/user"
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
			cErr := r.data.cache.Get(ctx, key, get)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				get, cErr = r.data.db.AccessLog.Get(ctx, logId)
			}
			return get, cErr
		})
	case v1.View_WITH_EDGE_IDS:
		// key: access_log_cache_key_get_access_log_id_edge_ids:logId
		key = r.cacheKey(strconv.FormatInt(logId, 10), append(r.ck["Get"], "edge_ids")...)
		res, err, _ = r.sg.Do(key, func() (any, error) {
			get := &ent.AccessLog{}
			// get cache
			cErr := r.data.cache.Get(ctx, key, get)
			if cErr != nil && errors.Is(cErr, cache.ErrCacheMiss) { // cache miss
				// get from db
				get, cErr = r.data.db.AccessLog.Query().
					Where(accesslog.ID(logId)).
					WithOwnerUser(func(query *ent.UserQuery) {
						query.Select(user.FieldID)
						query.Select(user.FieldToken)
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
			Value: res.(*ent.AccessLog),
			TTL:   r.data.conf.Cache.Ttl.AsDuration(),
		}); err != nil {
			r.log.Errorf("cache error: %v", err)
		}
		return toAccessLog(res.(*ent.AccessLog))
	case ent.IsNotFound(err): // db miss
		return nil, v1.ErrorNotFound("access log not found: %v", err)
	default: // error
		return nil, v1.ErrorUnknown("unknown error: %v", err)
	}
}

func (r *accessLogRepo) createBuilder(log *biz.AccessLog) *ent.AccessLogCreate {
	m := r.data.db.AccessLog.Create()
	m.SetIP(log.IP)
	m.SetCountryName(log.CountryName)
	m.SetCountryCode(log.CountryCode)
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
