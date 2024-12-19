package repo

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/suryasaputra2016/book-rental/entity"
)

type UserRepoMock struct {
	Mock mock.Mock
}

func (m *UserRepoMock) FindUserByID(id int) (*entity.User, error) {
	res := m.Mock.Called(id)

	if res.Get(0) == nil {
		return nil, errors.New("error")
	}

	product := res.Get(0).(entity.User)
	return &product, nil
}

// find user with email
func (ur *UserRepoMock) FindUserByEmail(email string) (*entity.User, error) {
	return nil, nil
}

// add user to database
func (ur *UserRepoMock) AddUser(userPtr *entity.User) error {
	return nil
}

func (ur *UserRepoMock) EditUser(userPtr *entity.User) error {
	return nil
}
