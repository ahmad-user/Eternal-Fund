package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"eternal-fund/mocking"
	"eternal-fund/model"
	"eternal-fund/model/dto"
)

type UserControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	userUseCase    *mocking.UserUseCaseMock
	authMiddleware *mocking.AuthMiddlewareMock
}

func (suite *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	suite.userUseCase = new(mocking.UserUseCaseMock)
	suite.authMiddleware = new(mocking.AuthMiddlewareMock)
	rg := r.Group("/api/v1")
	NewUserController(suite.userUseCase, rg, suite.authMiddleware).Routing()
	suite.router = r
}

func (suite *UserControllerTestSuite) TestListHandler_Success() {
	suite.userUseCase.On("FindAll", 1, 10).Return([]model.User{{ID: 1, Name: "User 1"}}, dto.Paging{TotalRows: 1, Size: 10, Page: 1}, nil)

	req, _ := http.NewRequest("GET", "/api/v1/users?page=1&size=10", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestGetByIdHandler_Success() {
	suite.userUseCase.On("FindById", 1).Return(model.User{ID: 1, Name: "User 1"}, nil)

	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestRegisterHandler_Success() {
	user := model.User{
		Name:       "John Doe",
		Email:      "john.doe@example.com",
		PasswordHash:   "password",
		Occupation: "Developer",
	}
	inputJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	suite.userUseCase.On("RegisterUser", mock.Anything).Return(user, nil)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUserHandler_Success() {
	user := model.User{
		ID:         1,
		Name:       "John Doe",
		Email:      "john.doe@example.com",
		PasswordHash:   "password",
		Occupation: "Developer",
	}
	inputJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	suite.userUseCase.On("UpdateUser", 1, mock.Anything).Return(user, nil)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestSaveAvatarHandler_Success() {
	file, _ := ioutil.ReadFile("test_avatar.jpg")
	req, _ := http.NewRequest("POST", "/api/v1/users/1/avatar", bytes.NewReader(file))
	req.Header.Set("Content-Type", "multipart/form-data")
	avatarFileName := "test_avatar.jpg"
	suite.userUseCase.On("SaveAvatar", 1, mock.Anything).Return(model.User{ID: 1, Name: "User 1", Email: "john.doe@example.com", AvatarFileName: &avatarFileName, PasswordHash: "password", Occupation: "Developer", Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestIsEmailAvailableHandler_Success() {
	checkEmailInput := model.CheckEmailInput{
		Email: "john.doe@example.com",
	}
	inputJSON, _ := json.Marshal(checkEmailInput)
	req, _ := http.NewRequest("POST", "/api/v1/users/check-email", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	suite.userUseCase.On("IsEmailAvailable", checkEmailInput).Return(true, nil)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *UserControllerTestSuite) TestListHandler_Error() {
	suite.userUseCase.On("FindAll", 1, 10).Return(nil, dto.Paging{}, errors.New("error"))

	req, _ := http.NewRequest("GET", "/api/v1/users?page=1&size=10", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
