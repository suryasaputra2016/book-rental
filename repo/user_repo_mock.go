package repo

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/suryasaputra2016/book-rental/entity"
)

// define mock user repository
type UserRepoMock struct {
	Mock mock.Mock
}

func (m *UserRepoMock) FindUserByID(id int) (*entity.User, error) {
	res := m.Mock.Called(id)

	if res.Get(0) == nil {
		return nil, errors.New("user id not found")
	}

	user := res.Get(0).(entity.User)
	return &user, nil
}

// FindUserByEmail gets user with email mock
func (m *UserRepoMock) FindUserByEmail(email string) (*entity.User, error) {
	res := m.Mock.Called(email)

	if res.Get(0) == nil {
		return nil, errors.New("user email not found")
	}

	user := res.Get(0).(entity.User)
	return &user, nil
}

// AddUser inserts user to database mock
func (m *UserRepoMock) AddUser(userPtr *entity.User) error {
	res := m.Mock.Called(userPtr)

	if res.Get(0) == nil {
		return errors.New("adding user failed")
	}

	return nil
}

// EditUser updates user in database mock
func (m *UserRepoMock) EditUser(userPtr *entity.User) error {
	res := m.Mock.Called(userPtr)

	if res.Get(0) == nil {
		return errors.New("edit user failed")
	}

	return nil
}
