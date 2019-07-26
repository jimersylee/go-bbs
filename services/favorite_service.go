package services

import (
	"errors"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
)

var FavoriteService = newFavoriteService()

func newFavoriteService() *favoriteService {
	return &favoriteService{}
}

type favoriteService struct {
}

func (this *favoriteService) Get(id int64) *model.Favorite {
	return repositories.FavoriteRepository.Get(simple.GetDB(), id)
}

func (this *favoriteService) Take(where ...interface{}) *model.Favorite {
	return repositories.FavoriteRepository.Take(simple.GetDB(), where...)
}

func (this *favoriteService) QueryCnd(cnd *simple.QueryCnd) (list []model.Favorite, err error) {
	return repositories.FavoriteRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *favoriteService) Query(queries *simple.ParamQueries) (list []model.Favorite, paging *simple.Paging) {
	return repositories.FavoriteRepository.Query(simple.GetDB(), queries)
}

func (this *favoriteService) Create(t *model.Favorite) error {
	return repositories.FavoriteRepository.Create(simple.GetDB(), t)
}

func (this *favoriteService) Update(t *model.Favorite) error {
	return repositories.FavoriteRepository.Update(simple.GetDB(), t)
}

func (this *favoriteService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.FavoriteRepository.Updates(simple.GetDB(), id, columns)
}

func (this *favoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.FavoriteRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *favoriteService) Delete(id int64) {
	repositories.FavoriteRepository.Delete(simple.GetDB(), id)
}

func (this *favoriteService) GetBy(entityType string, entityId int64) *model.Favorite {
	return repositories.FavoriteRepository.Take(simple.GetDB(), "entity_type = ? and entity_id = ?", entityType, entityId)
}

// 收藏文章
func (this *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := repositories.ArticleRepository.Get(simple.GetDB(), articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		return errors.New("收藏的文章不存在")
	}
	temp := this.GetBy(model.EntityTypeArticle, articleId)
	if temp != nil { // 已经收藏
		return nil
	}
	return repositories.FavoriteRepository.Create(simple.GetDB(), &model.Favorite{
		UserId:     userId,
		EntityType: model.EntityTypeArticle,
		EntityId:   articleId,
		CreateTime: simple.NowTimestamp(),
	})
}

// 收藏主题
func (this *favoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := repositories.TopicRepository.Get(simple.GetDB(), topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return errors.New("收藏的话题不存在")
	}
	temp := this.GetBy(model.EntityTypeTopic, topicId)
	if temp != nil { // 已经收藏
		return nil
	}
	return repositories.FavoriteRepository.Create(simple.GetDB(), &model.Favorite{
		UserId:     userId,
		EntityType: model.EntityTypeTopic,
		EntityId:   topicId,
		CreateTime: simple.NowTimestamp(),
	})
}
