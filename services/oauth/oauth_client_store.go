package oauth

import (
	"errors"
	"github.com/jimersylee/go-bbs/model"
	"github.com/jimersylee/go-bbs/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
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
