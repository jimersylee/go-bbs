package admin

import (
	"encoding/json"
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris"
	"strconv"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (this *SysConfigController) GetBy(id int64) *simple.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *SysConfigController) AnyList() *simple.JsonResult {
	list, paging := services.SysConfigService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *SysConfigController) GetAll() *simple.JsonResult {
	list := services.SysConfigService.GetAll()
	return simple.JsonData(list)
}

func (this *SysConfigController) PostSave() *simple.JsonResult {
	config := this.Ctx.FormValue("config")
	data := make(map[string]string)
	err := json.Unmarshal([]byte(config), &data)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	err = services.SysConfigService.SetAll(data)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}
