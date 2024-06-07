package mocking

import (
	"eternal-fund/model"
	"eternal-fund/model/dto"

	"github.com/stretchr/testify/mock"
)

// UserUseCaseMock adalah implementasi mock dari UserUseCase.
type UserUseCaseMock struct {
	mock.Mock
}

// Save memock metode Save dari UserUseCase.
func (m *UserUseCaseMock) Save(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

// RegisterUser memock metode RegisterUser dari UserUseCase.
func (m *UserUseCaseMock) RegisterUser(input model.RegisterUserInput) (model.User, error) {
	args := m.Called(input)
	return args.Get(0).(model.User), args.Error(1)
}

// Update memock metode Update dari UserUseCase.
func (m *UserUseCaseMock) Update(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

// SaveAvatar memock metode SaveAvatar dari UserUseCase.
func (m *UserUseCaseMock) SaveAvatar(userId int, fileLocation string) (model.User, error) {
	args := m.Called(userId, fileLocation)
	return args.Get(0).(model.User), args.Error(1)
}

// IsEmailAvailable memock metode IsEmailAvailable dari UserUseCase.
func (m *UserUseCaseMock) IsEmailAvailable(input model.CheckEmailInput) (bool, error) {
	args := m.Called(input)
	return args.Bool(0), args.Error(1)
}

// FindAll memock metode FindAll dari UserUseCase.
func (m *UserUseCaseMock) FindAll(page int, size int) ([]model.User, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.User), args.Get(1).(dto.Paging), args.Error(2)
}

// FindById memock metode FindById dari UserUseCase.
func (m *UserUseCaseMock) FindById(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

// FindByEmail memock metode FindByEmail dari UserUseCase.
func (m *UserUseCaseMock) FindByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserUseCaseMock) UpdateUser(id int, input model.User) (model.User, error) {
	args := m.Called(id, input)
	return args.Get(0).(model.User), args.Error(1)
}

// NewUserUseCaseMock membuat instance baru dari UserUseCaseMock.
func NewUserUseCaseMock() *UserUseCaseMock {
	return &UserUseCaseMock{}
}
