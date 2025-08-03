package repository

import (
	"github.com/yuheiLin/go-http-server/customerror"
	"github.com/yuheiLin/go-http-server/model"
	"sync"
)

type Repository interface {
	GetUser(userID string) (*model.User, error)
	CreateUser(userID string, password string) (*model.User, error)
	VerifyUser(userID string, password string) error
}

type repositoryImpl struct {
	userData *sync.Map
}

type userVal struct {
	UserID   string
	Nickname string
	Password string
	Comment  string
}

func New() Repository {
	// Initialize the repository with an empty sync.Map for user data
	m := sync.Map{}

	firstUser := userVal{
		UserID:   "TaroYamada",
		Password: "PaSSwd4TY",
		Nickname: "たろー",
		Comment:  "僕は元気です",
	}

	// insert a test data
	m.Store(firstUser.UserID, firstUser)

	return &repositoryImpl{
		userData: &m,
	}
}

func (r *repositoryImpl) GetUser(userID string) (*model.User, error) {
	u, ok := r.userData.Load(userID)
	if !ok {
		return nil, customerror.ErrUserNotFound
	}
	ud := u.(userVal)
	um := &model.User{
		ID:       ud.UserID,
		Nickname: ud.Nickname,
		Comment:  ud.Comment,
	}
	return um, nil
}

func (r *repositoryImpl) CreateUser(userID string, password string) (*model.User, error) {
	_, ok := r.userData.Load(userID)
	if ok {
		return nil, customerror.ErrUserAlreadyExists
	}

	user := userVal{
		UserID:   userID,
		Password: password,
	}

	r.userData.Store(user.UserID, user)
	return &model.User{
		ID:       user.UserID,
		Nickname: user.Nickname,
		Comment:  user.Comment,
	}, nil
}

func (r *repositoryImpl) VerifyUser(userID string, password string) error {
	u, ok := r.userData.Load(userID)
	if !ok {
		return customerror.ErrUserNotFound
	}
	ud := u.(userVal)

	if ud.Password != password {
		return customerror.ErrUserNotFound
	}

	return nil
}

//func (r *repositoryImpl) DeleteUser(userID string) error {
//	r.userData.Delete(userID)
//	return nil
//}
