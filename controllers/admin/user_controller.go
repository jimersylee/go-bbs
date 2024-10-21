package admin

import (
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/utils"
	"github.com/kataras/iris/v12/context"
	"strconv"
	"strings"

	"github.com/jimersylee/go-bbs/services"
	"github.com/jimersylee/go-bbs/utils/simple"
)

type UserController struct {
	Ctx *context.Context
}

func (this *UserController) GetBy(id int64) *simple.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(this.buildUserItem(t))
}

func (this *UserController) AnyList() *simple.JsonResult {
	list, paging := services.UserService.Query(simple.NewParamQueries(this.Ctx).EqAuto("id").EqAuto("username").PageAuto().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, this.buildUserItem(&user))
	}
	return simple.JsonData(&simple.PageResult{Results: itemList, Page: paging})
}

func (this *UserController) PostCreate() *simple.JsonResult {
	username := simple.FormValue(this.Ctx, "username")
	email := simple.FormValue(this.Ctx, "email")
	password := simple.FormValue(this.Ctx, "password")
	nickname := simple.FormValue(this.Ctx, "nickname")

	user, err := services.UserService.SignUp(username, email, password, password, nickname, "")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(this.buildUserItem(user))
}

func (this *UserController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.UserService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	username := simple.FormValue(this.Ctx, "username")
	password := simple.FormValue(this.Ctx, "password")
	nickname := simple.FormValue(this.Ctx, "nickname")
	email := simple.FormValue(this.Ctx, "email")
	roles := simple.FormValueStringArray(this.Ctx, "roles")
	status := simple.FormValueIntDefault(this.Ctx, "status", -1)

	if len(username) > 0 {
		t.Username = username
	}
	if len(password) > 0 {
		t.Password = simple.EncodePassword(t.Password)
	}
	if len(nickname) > 0 {
		t.Nickname = nickname
	}
	if len(email) > 0 {
		t.Email = email
	}
	if status != -1 {
		t.Status = status
	}

	t.Roles = strings.Join(roles, ",")

	err = services.UserService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *UserController) buildUserItem(user *model.User) map[string]interface{} {
	return simple.NewRspBuilder(user).Put("roles", utils.GetUserRoles(user.Roles)).Build()
}
