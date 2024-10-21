package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/jimersylee/go-bbs/controllers/render"
	"github.com/jimersylee/go-bbs/model"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetUserInfoByRequest(r *http.Request) *model.UserInfo {
	token, err := ServerInstance.Srv.ValidationBearerToken(r)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return GetUserInfoByToken(token)
}

func GetUserInfo(accessToken string) *model.UserInfo {
	if len(accessToken) == 0 {
		return nil
	}
	token, err := ServerInstance.Srv.Manager.LoadAccessToken(nil, accessToken)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return GetUserInfoByToken(token)
}

func GetUserInfoByToken(token oauth2.TokenInfo) *model.UserInfo {
	userId, err := strconv.ParseInt(token.GetUserID(), 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	if userId <= 0 {
		return nil
	}
	return render.BuildUserById(userId)
}
