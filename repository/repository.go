package repository

import (
	"errors"
	"github.com/yuheiLin/go-http-server/model"
	"sync"
)

type Repository interface {
	GetUser(userID string) (*model.User, error)
	CreateUser(userID string, userName string) error
	DeleteUser(userID string) error
}

type repositoryImpl struct {
	userData *sync.Map
}

func New() Repository {
	return &repositoryImpl{
		userData: &sync.Map{},
	}
}

type userVal struct {
	ID   string
	Name string
}

func (r *repositoryImpl) GetUser(userID string) (*model.User, error) {
	u, ok := r.userData.Load(userID)
	if !ok {
		return nil, errors.New("user not found")
	}
	ud := u.(userVal)
	um := &model.User{
		ID:   ud.ID,
		Name: ud.Name,
	}
	return um, nil
}

func (r *repositoryImpl) CreateUser(userID string, userName string) error {
	r.userData.Store(userID, userVal{
		ID:   userID,
		Name: userName,
	})
	return nil
}

func (r *repositoryImpl) DeleteUser(userID string) error {
	r.userData.Delete(userID)
	return nil
}
