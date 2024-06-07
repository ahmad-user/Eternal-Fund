package mocking

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) FindById(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) FindAll(page int, size int) ([]model.User, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.User), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *UserRepoMock) FindByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) Save(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) Update(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) SaveAvatar(userId int, fileLocation string) (model.User, error) {
	args := m.Called(userId, fileLocation)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserRepoMock) IsEmailAvailable(input model.CheckEmailInput) (bool, error) {
	args := m.Called(input)
	return args.Bool(0), args.Error(1)
}

func NewUserRepoMock(db *sql.DB) *UserRepoMock {
	return &UserRepoMock{}
}