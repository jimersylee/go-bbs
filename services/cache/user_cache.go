package cache

import (
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"time"

	"github.com/goburrow/cache"
	"github.com/jimersylee/go-bbs/utils/simple"
)

type userCache struct {
	cache           cache.LoadingCache
	activeUserCache cache.LoadingCache // 活跃用户缓存
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserRepository.Get(simple.GetDB(), Key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		activeUserCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				dateFrom := simple.Timestamp(simple.WithTimeAsStartOfDay(time.Now()))
				rows, e := simple.GetDB().Raw("select user_id, count(*) c from t_article where create_time > ?"+
					" group by user_id order by c desc limit 20", dateFrom).Rows()
				// rows, e := simple.GetDB().Raw("select user_id, count(*) c from t_article " +
				// 	" group by user_id order by c desc limit 10").Rows()
				if e != nil {
					return
				}
				var userIds []int64
				for rows.Next() {
					var userId int64
					var c int
					err := rows.Scan(&userId, &c)
					if err != nil {
						continue
					}
					userIds = append(userIds, userId)
				}
				value = userIds
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (this *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	val, err := this.cache.Get(userId)
	if err != nil {
		return nil
	}
	return val.(*model.User)
}

func (this *userCache) Invalidate(userId int64) {
	this.cache.Invalidate(userId)
}

func (this *userCache) GetActiveUsers() []model.User {
	val, err := this.activeUserCache.Get("data")
	if err != nil {
		return nil
	}
	userIds := val.([]int64)
	if len(userIds) == 0 {
		return nil
	}
	var users []model.User
	for _, userId := range userIds {
		user := this.Get(userId)
		if user != nil {
			users = append(users, *user)
		}
	}
	return users
}
