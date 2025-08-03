package service

import (
	"github.com/yuheiLin/go-http-server/apiclient"
	"github.com/yuheiLin/go-http-server/model"
	"github.com/yuheiLin/go-http-server/repository"
)

type Service interface {
	GetUser(userID string) (*model.User, error)
	CreateUser(userID string, userName string) error
	DeleteUser(userID string) error
}
type serviceImpl struct {
	repository repository.Repository
	apiclient  apiclient.ApiClient
}

func New(repository repository.Repository, apiclient apiclient.ApiClient) Service {
	return &serviceImpl{
		repository: repository,
		apiclient:  apiclient,
	}
}

func (service *serviceImpl) GetUser(userID string) (*model.User, error) {
	return service.repository.GetUser(userID)
}

func (service *serviceImpl) CreateUser(userID string, userName string) error {
	return service.repository.CreateUser(userID, userName)
}

func (service *serviceImpl) DeleteUser(userID string) error {
	return service.repository.DeleteUser(userID)
}
