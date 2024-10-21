package admin

import (
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris/v12/context"
	"strconv"
)

type TopicTagController struct {
	Ctx *context.Context
}

func (this *TopicTagController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TopicTagController) AnyList() *simple.JsonResult {
	list, paging := services.TopicTagService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TopicTagController) PostCreate() *simple.JsonResult {
	t := &model.TopicTag{}
	this.Ctx.ReadForm(t)

	err := services.TopicTagService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TopicTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.TopicTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.TopicTagService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
