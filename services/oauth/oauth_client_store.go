package oauth

import (
	"errors"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
)

type clientStore struct {
}

func newClientStore() *clientStore {
	return &clientStore{}
}

func (s *clientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	oauthClient := repositories.OauthClientRepository.GetByClientId(simple.GetDB(), id)
	if oauthClient == nil {
		return nil, errors.New("Client not found:" + id)
	}
	if oauthClient.Status == model.OauthClientStatusDisabled {
		return nil, errors.New("Client disabled:" + id)
	}
	return &models.Client{
		ID:     oauthClient.ClientId,
		Secret: oauthClient.ClientSecret,
		Domain: oauthClient.Domain,
	}, nil
}
