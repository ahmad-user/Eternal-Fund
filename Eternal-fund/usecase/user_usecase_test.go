package usecase

import (
	"database/sql"
	"errors"
	"eternal-fund/mocking"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UsersUseCaseTestSuite struct {
	suite.Suite
	uuc          *userUseCase
	userRepoMock *mocking.UserUseCaseMock
}

func (suite *UsersUseCaseTestSuite) SetupTest() {
	suite.userRepoMock = new(mocking.UserUseCaseMock)
	suite.uuc = &userUseCase{repo: suite.userRepoMock}
}
func (suite *UsersUseCaseTestSuite) TestIsEmailAvailable_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo
	email := "test@example.com"
	input := model.CheckEmailInput{Email: email}
	mockUser := model.User{Email: email /* other fields */}
	mockRepo.On("FindByEmail", email).Return(mockUser, nil)
	available, err := suite.uuc.IsEmailAvailable(input)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.False(suite.T(), available, "Expected email to not be available")
}

func (suite *UsersUseCaseTestSuite) TestSaveAvatar_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo
	userId := 1
	fileLocation := "/path/to/avatar.jpg"
	expectedUser := model.User{
		ID:             userId,
		AvatarFileName: &fileLocation,
	}
	mockRepo.On("SaveAvatar", userId, fileLocation).Return(expectedUser, nil)
	user, err := suite.uuc.SaveAvatar(userId, fileLocation)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser, user, "Expected saved user object to match")
}

func (suite *UsersUseCaseTestSuite) TestFindAll_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo
	page := 1
	size := 10
	expectedUsers := []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	expectedPaging := dto.Paging{
		Page: page,
		Size: size,
	}
	mockRepo.On("FindAll", page, size).Return(expectedUsers, expectedPaging, nil)
	users, paging, err := suite.uuc.FindAll(page, size)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUsers, users, "Expected returned users to match")
	assert.Equal(suite.T(), expectedPaging, paging, "Expected returned paging info to match")
}

func (suite *UsersUseCaseTestSuite) TestFindById_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo
	userID := 1
	expectedUser := model.User{
		ID:   userID,
		Name: "Alice",
	}
	mockRepo.On("FindById", userID).Return(expectedUser, nil)
	user, err := suite.uuc.FindById(userID)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser, user, "Expected returned user to match")
}

func (suite *UsersUseCaseTestSuite) TestFindByEmail_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo
	email := "test@example.com"
	expectedUser := model.User{
		ID:    1,
		Email: email,
	}
	mockRepo.On("FindByEmail", email).Return(expectedUser, nil)
	user, err := suite.uuc.FindByEmail(email)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser, user, "Expected returned user to match")
}

func (suite *UsersUseCaseTestSuite) TestRegister_Success() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo

	// Define test input
	input := model.RegisterUserInput{
		Name:       "John Doe",
		Occupation: "Software Developer",
		Email:      "john.doe@example.com",
		Password:   "securepassword",
	}
	expectedUser := model.User{
		ID:           1,
		Name:         input.Name,
		Occupation:   input.Occupation,
		Email:        input.Email,
		PasswordHash: "hashed_password",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	mockRepo.On("Save", mock.AnythingOfType("model.User")).Return(expectedUser, nil)
	registeredUser, err := suite.uuc.RegisterUser(input)

	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser.Name, registeredUser.Name, "Expected registered user's name to match")
	assert.Equal(suite.T(), expectedUser.Email, registeredUser.Email, "Expected registered user's email to match")
}

func (suite *UsersUseCaseTestSuite) TestRegisterUser_FailedBySaveError() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo

	input := model.RegisterUserInput{
		Name:       "John Doe",
		Occupation: "Software Developer",
		Email:      "john.doe@example.com",
		Password:   "securepassword",
	}

	mockRepo.On("Save", mock.AnythingOfType("model.User")).Return(model.User{}, errors.New("failed to save user"))
	_, err := suite.uuc.RegisterUser(input)

	assert.Error(suite.T(), err, "Expected an error")
	assert.Equal(suite.T(), "failed to save user", err.Error(), "Expected error message to match")
}

func (suite *UsersUseCaseTestSuite) TestUpdateUser_FailedByNotFound() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo

	userId := 1
	input := model.User{
		Name:       "Updated Name",
		Occupation: "Updated Occupation",
		Email:      "updated@example.com",
	}

	mockRepo.On("FindById", userId).Return(model.User{}, sql.ErrNoRows)
	_, err := suite.uuc.UpdateUser(userId, input)

	assert.Error(suite.T(), err, "Expected an error")
	assert.True(suite.T(), errors.Is(err, sql.ErrNoRows), "Expected sql.ErrNoRows error")
}

func (suite *UsersUseCaseTestSuite) TestIsEmailAvailable_FailedByUnexpectedError() {
	mockRepo := &mocking.UserUseCaseMock{}
	suite.uuc.repo = mockRepo

	email := "test@example.com"
	input := model.CheckEmailInput{Email: email}

	mockRepo.On("FindByEmail", email).Return(model.User{}, errors.New("unexpected error"))
	_, err := suite.uuc.IsEmailAvailable(input)

	assert.Error(suite.T(), err, "Expected an error")
	assert.Equal(suite.T(), "unexpected error", err.Error(), "Expected error message to match")
}

func TestUsersUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UsersUseCaseTestSuite))
}
