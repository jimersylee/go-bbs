package admin

import (
	"github.com/kataras/iris/v12/context"
	"strconv"

	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
)

type UserArticleTagController struct {
	Ctx *context.Context
}

func (this *UserArticleTagController) GetBy(id int64) *simple.JsonResult {
	t := services.UserArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *UserArticleTagController) AnyList() *simple.JsonResult {
	list, paging := services.UserArticleTagService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *UserArticleTagController) PostCreate() *simple.JsonResult {
	t := &model.UserArticleTag{}
	this.Ctx.ReadForm(t)

	err := services.UserArticleTagService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *UserArticleTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.UserArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.UserArticleTagService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
