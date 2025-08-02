package service

import (
	"github.com/yuheiLin/go-http-server/apiclient"
	"github.com/yuheiLin/go-http-server/repository"
)

type Service interface {
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
