package oauth

import (
	"errors"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/json-iterator/go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// TokenStore实现
type tokenStore struct {
}

// 存储oauth token
func newTokenStore() *tokenStore {
	return &tokenStore{}
}

func (s *tokenStore) Create(info oauth2.TokenInfo) error {
	buf, _ := jsoniter.Marshal(info)
	item := &model.OauthToken{
		Data: string(buf),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.AccessToken = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()
		if refresh := info.GetRefresh(); refresh != "" {
			item.RefreshToken = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
		}
	}
	return repositories.OauthTokenRepository.Create(simple.GetDB(), item)
}

func (s *tokenStore) RemoveByCode(code string) error {
	repositories.OauthTokenRepository.RemoveByCode(simple.GetDB(), code)
	return nil
}

func (s *tokenStore) RemoveByAccess(access string) error {
	repositories.OauthTokenRepository.RemoveByAccessToken(simple.GetDB(), access)
	return nil
}

func (s *tokenStore) RemoveByRefresh(refresh string) error {
	repositories.OauthTokenRepository.RemoveByRefreshToken(simple.GetDB(), refresh)
	return nil
}

func (s *tokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	if len(code) == 0 {
		return nil, nil
	}
	oauthToken := repositories.OauthTokenRepository.GetByCode(simple.GetDB(), code)
	if oauthToken == nil {
		return nil, errors.New("invalidate code")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *tokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	if len(access) == 0 {
		return nil, nil
	}
	oauthToken := repositories.OauthTokenRepository.GetByAccessToken(simple.GetDB(), access)
	if oauthToken == nil {
		return nil, errors.New("invalidate access token")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *tokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	if len(refresh) == 0 {
		return nil, nil
	}
	oauthToken := repositories.OauthTokenRepository.GetByRefreshToken(simple.GetDB(), refresh)
	if oauthToken == nil {
		return nil, errors.New("invalidate refresh token")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *tokenStore) toTokenInfo(data string) (oauth2.TokenInfo, error) {
	var tm models.Token
	err := jsoniter.Unmarshal([]byte(data), &tm)
	if err != nil {
		return nil, err
	}
	return &tm, nil
}
