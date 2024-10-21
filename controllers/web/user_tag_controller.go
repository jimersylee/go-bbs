package web

import (
	"github.com/jimersylee/go-bbs/utils/session"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"

	"github.com/jimersylee/go-bbs/controllers/render"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/services/cache"
)

type UserTagController struct {
	Ctx *context.Context
}

// 标签列表页面
func (this *UserTagController) GetList() mvc.View {
	user := session.GetCurrentUser(this.Ctx)
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	list, paging := services.UserArticleTagService.Query(simple.NewParamQueries(this.Ctx).Eq("user_id", user.Id).Page(page, 20).Desc("id"))

	var tags []model.Tag
	if len(list) > 0 {
		for _, v := range list {
			tags = append(tags, *cache.TagCache.Get(v.TagId))
		}
	}
	return mvc.View{
		Name: "user/tag/list.html",
		Data: iris.Map{
			"User":             render.BuildUser(user),
			"Tags":             render.BuildTags(tags),
			"Page":             paging,
			model.TplSiteTitle: user.Nickname + " - 标签列表",
		},
	}
}
