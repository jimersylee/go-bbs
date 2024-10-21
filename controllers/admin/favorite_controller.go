package admin

import (
	"github.com/kataras/iris/v12/context"
	"strconv"

	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
)

type FavoriteController struct {
	Ctx *context.Context
}

func (this *FavoriteController) GetBy(id int64) *simple.JsonResult {
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) AnyList() *simple.JsonResult {
	list, paging := services.FavoriteService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *FavoriteController) PostCreate() *simple.JsonResult {
	t := &model.Favorite{}
	this.Ctx.ReadForm(t)

	err := services.FavoriteService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.FavoriteService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
