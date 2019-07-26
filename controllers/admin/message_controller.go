package admin

import (
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris"
	"strconv"
)

type MessageController struct {
	Ctx iris.Context
}

func (this *MessageController) GetBy(id int64) *simple.JsonResult {
	t := services.MessageService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *MessageController) AnyList() *simple.JsonResult {
	list, paging := services.MessageService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *MessageController) PostCreate() *simple.JsonResult {
	t := &model.Message{}
	this.Ctx.ReadForm(t)

	err := services.MessageService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *MessageController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.MessageService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.MessageService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
