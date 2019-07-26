package web

import (
	context2 "context"
	"github.com/jimersylee/go-bbs/services/oauth"
	"github.com/jimersylee/go-bbs/utils/config"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"golang.org/x/oauth2"
)

const oauthState = "mlog"

type OauthClientController struct {
	Ctx context.Context
}

// 跳转到授权页
func (c *OauthClientController) Any() {
	authCodeUrl := oauth.GetOauthConfig().AuthCodeURL(oauthState)
	c.Ctx.Redirect(authCodeUrl, iris.StatusSeeOther)
}

// 回调地址
func (c *OauthClientController) AnyCallback() {
	state := c.Ctx.FormValue("state")
	if state != oauthState {
		_, _ = c.Ctx.JSON(simple.ErrorMsg("State invalid"))
		return
	}
	code := c.Ctx.FormValue("code")
	token, err := oauth.GetOauthConfig().Exchange(context2.TODO(), code)
	if err != nil {
		_, _ = c.Ctx.JSON(simple.ErrorMsg("Code exchange failed with " + err.Error()))
		return
	}
	_, _ = c.Ctx.HTML(oauth.GetSuccessHtml(token, config.Conf.OauthClient.ClientSuccessUrl))
}

// 通过refreshToken重新换取accessToken
func (c *OauthClientController) AnyRefresh() *simple.JsonResult {
	refreshToken := c.Ctx.FormValue("refreshToken")
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := oauth.GetOauthConfig().TokenSource(context2.TODO(), token)
	newToken, err := ts.Token()
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(newToken)
}
