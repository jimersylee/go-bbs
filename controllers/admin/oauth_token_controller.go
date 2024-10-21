package admin

import (
	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris/v12/context"
	"strconv"
)

type OauthTokenController struct {
	Ctx *context.Context
}

func (this *OauthTokenController) GetBy(id int64) *simple.JsonResult {
	t := services.OauthTokenService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *OauthTokenController) AnyList() *simple.JsonResult {
	list, paging := services.OauthTokenService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}
