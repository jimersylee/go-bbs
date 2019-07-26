package web

import (
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/services/cache"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/jimersylee/go-bbs/controllers/render"
	"github.com/jimersylee/go-bbs/model"
)

type IndexController struct {
	Ctx iris.Context
}

func (this *IndexController) Any() mvc.View {
	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	activeTags := cache.TagCache.GetActiveTags()

	// 存在缓存从缓存里面取
	articles := cache.ArticleCache.GetIndexList()

	topics, _ := services.TopicService.QueryCnd(simple.NewQueryCnd("status = ?", model.TopicStatusOk).Order("id desc").Size(10))

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{
			model.TplSiteDescription: cache.SysConfigCache.GetValue(model.SysConfigSiteDescription),
			model.TplSiteKeywords:    cache.SysConfigCache.GetValue(model.SysConfigSiteKeywords),
			"Categories":             categories,
			"Articles":               render.BuildArticles(articles),
			"Topics":                 render.BuildTopics(topics),
			"ActiveUsers":            render.BuildUsers(activeUsers),
			"ActiveTags":             render.BuildTags(activeTags),
		},
	}
}

func (this *IndexController) GetAbout() mvc.View {
	return mvc.View{
		Name: "about.html",
	}
}
