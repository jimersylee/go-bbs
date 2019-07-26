package cache

import (
	"github.com/goburrow/cache"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/sirupsen/logrus"
	"time"
)

type categoryCache struct {
	cache              cache.LoadingCache
	allCategoriesCache cache.LoadingCache
}

var CategoryCache = newCategoryCache()

func newCategoryCache() *categoryCache {
	return &categoryCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.CategoryRepository.Get(simple.GetDB(), Key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		allCategoriesCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				list, e := repositories.CategoryRepository.GetCategories()
				if e != nil {
					return
				}

				var categories []model.CategoryResponse
				for _, cat := range list {
					categories = append(categories, model.CategoryResponse{
						CategoryId:   cat.Id,
						CategoryName: cat.Name,
					})
				}
				value = categories
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (this *categoryCache) Get(categoryId int64) *model.Category {
	val, err := this.cache.Get(categoryId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if val != nil {
		return val.(*model.Category)
	}
	return nil
}

func (this *categoryCache) Invalidate(categoryId int64) {
	this.cache.Invalidate(categoryId)
}

func (this *categoryCache) GetAllCategories() []model.CategoryResponse {
	val, err := this.allCategoriesCache.Get("data")
	if err != nil {
		return nil
	}
	return val.([]model.CategoryResponse)
}

func (this *categoryCache) InvalidateAll() {
	this.allCategoriesCache.InvalidateAll()
}
