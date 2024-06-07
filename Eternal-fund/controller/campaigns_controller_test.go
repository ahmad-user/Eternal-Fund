package controller

import (
	"bytes"
	"encoding/json"
	"eternal-fund/mocking"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CampaignsControllerTestSuite struct {
	suite.Suite
	rg  *gin.RouterGroup
	aum *mocking.CampaignsUseCaseMock
	amm *mocking.AuthMiddlewareMock
}

func (suite *CampaignsControllerTestSuite) SetupTest() {
	suite.aum = new(mocking.CampaignsUseCaseMock)
	suite.amm = new(mocking.AuthMiddlewareMock)
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rg := r.Group("api/v1")

	rg.Use(suite.amm.CheckToken("user"))
	suite.rg = rg
}

func (suite *CampaignsControllerTestSuite) TestAllCampaigns_success() {
	var mockAuthor = []model.Campaigns{
		{
			ID:                1,
			User_id:           101,
			Name:              "Campaign 1",
			Short_description: "Short description 1",
			Description:       "Description 1",
			Perks:             "Perks 1",
			Backer_count:      10,
			Goal_amount:       1000,
			Current_amount:    500,
			Slug:              "campaign-1",
			Created_at:        time.Now(),
			Updated_at:        time.Now(),
		},
	}
	moackPaging := dto.Paging{
		TotalRows: 1,
		Size:      5,
		Page:      1,
	}
	suite.aum.On("FindAllCampaigns", 1, 5).Return(mockAuthor, moackPaging, nil)

	authorController := NewCampaignsController(suite.aum, suite.rg, suite.amm)
	authorController.Routing()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/campaigns?page=1&size=5", nil)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request
	authorController.getCampaignsHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *CampaignsControllerTestSuite) TestListHandler_fail() {

	suite.aum.On("FindAllCampaigns", 1, 5).Return([]model.Campaigns{}, dto.Paging{}, fmt.Errorf("error"))

	authorController := NewCampaignsController(suite.aum, suite.rg, suite.amm)
	authorController.Routing()

	request, err := http.NewRequest(http.MethodGet, "api/v1/authors?page=1&size=5", nil)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request
	authorController.getCampaignsHandler(ctx)
	assert.Equal(suite.T(), http.StatusInternalServerError, record.Code)
}

func (suite *CampaignsControllerTestSuite) TestDeleteCampaigns_success() {
	mockCampaignID := 1
	suite.aum.On("DeleteCampaigns", mockCampaignID).Return(nil)
	campaignController := NewCampaignsController(suite.aum, suite.rg, suite.amm)
	campaignController.Routing()
	request, err := http.NewRequest(http.MethodDelete, "/api/v1/campaigns/1", nil)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	campaignController.deleteCampaignHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
	var response dto.SingleResponse
	err = json.Unmarshal(record.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Status.Code)
	assert.Equal(suite.T(), "Campaign deleted successfully", response.Status.Message)
	assert.Nil(suite.T(), response.Data)
}

func (suite *CampaignsControllerTestSuite) TestByIdCampaigns_success() {
	mockCampaign := model.Campaigns{
		ID:                1,
		User_id:           101,
		Name:              "Campaign 1",
		Short_description: "Short description 1",
		Description:       "Description 1",
		Perks:             "Perks 1",
		Backer_count:      10,
		Goal_amount:       1000,
		Current_amount:    500,
		Slug:              "campaign-1",
		Created_at:        time.Now(),
		Updated_at:        time.Now(),
	}
	suite.aum.On("FindByIdCampaigns", mock.Anything).Return(mockCampaign, nil)
	campaignController := NewCampaignsController(suite.aum, suite.rg, suite.amm)
	campaignController.Routing()
	request, err := http.NewRequest(http.MethodGet, "/api/v1/campaigns/1", nil)
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request
	ctx.Params = gin.Params{{
		Key: "id", Value: "1"},
	}
	campaignController.getCampaignByIdHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}
func (suite *CampaignsControllerTestSuite) TestUpdateCampaignHandler_success() {
	mockUpdatedCampaign := model.Campaigns{
		ID:                1,
		User_id:           101,
		Name:              "Updated Campaign 1",
		Short_description: "Updated Short description 1",
		Description:       "Updated Description 1",
		Perks:             "Updated Perks 1",
		Backer_count:      15,
		Goal_amount:       1200,
		Current_amount:    700,
		Slug:              "updated-campaign-1",
		Updated_at:        time.Now(),
		User: model.User{
			ID:           101,
			Name:         "John Doe",
			Occupation:   "Developer",
			Email:        "john.doe@example.com",
			PasswordHash: "",
			Role:         "user",
		},
	}
	suite.aum.On("UpdateCampaigns", mock.AnythingOfType("int")).Return(mockUpdatedCampaign, nil)
	campaignController := NewCampaignsController(suite.aum, suite.rg, suite.amm)
	campaignController.Routing()
	updatePayload := []byte(`{
		"name": "Updated Campaign 1",
		"short_description": "Updated Short description 1",
		"description": "Updated Description 1",
		"perks": "Updated Perks 1",
		"backer_count": 15,
		"goal_amount": 1200,
		"current_amount": 700,
		"slug": "updated-campaign-1"
	}`)
	request, err := http.NewRequest(http.MethodPut, "/api/v1/campaigns/1", bytes.NewBuffer(updatePayload))
	assert.NoError(suite.T(), err)
	record := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(record)
	ctx.Request = request
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	campaignController.updateCampaignHandler(ctx)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
	var response dto.SingleResponse
	err = json.Unmarshal(record.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Status.Code)
	assert.Equal(suite.T(), "Campaign updated successfully", response.Status.Message)
	assert.Empty(suite.T(), mockUpdatedCampaign.User.PasswordHash, "PasswordHash should be empty after update")
	assert.True(suite.T(), mockUpdatedCampaign.User.CreatedAt.IsZero(), "CreatedAt should be time.Time{} after update")
	assert.True(suite.T(), mockUpdatedCampaign.User.UpdatedAt.IsZero(), "UpdatedAt should be time.Time{} after update")
}

func (suite *CampaignsControllerTestSuite) TestCreateCampaignHandler_success() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	mockInput := model.Campaigns{
		ID:                1,
		User_id:           101,
		Name:              "Updated Campaign 1",
		Short_description: "Updated Short description 1",
		Description:       "Updated Description 1",
		Perks:             "Updated Perks 1",
		Backer_count:      15,
		Goal_amount:       1200,
		Current_amount:    700,
		Slug:              "updated-campaign-1",
		Created_at:        time.Now(),
		User: model.User{
			ID:           101,
			Name:         "John Doe",
			Occupation:   "Developer",
			Email:        "john.doe@example.com",
			PasswordHash: "",
			Role:         "user",
		},
	}
	inputJSON, _ := json.Marshal(mockInput)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/api/campaigns", bytes.NewReader(inputJSON))
	ginContext.Request.Header.Set("Content-Type", "application/json")
	mockCampaign := model.Campaigns{
		ID:                1,
		User_id:           101,
		Name:              "Updated Campaign 1",
		Short_description: "Updated Short description 1",
		Description:       "Updated Description 1",
		Perks:             "Updated Perks 1",
		Backer_count:      15,
		Goal_amount:       1200,
		Current_amount:    700,
		Slug:              "updated-campaign-1",
		Updated_at:        time.Now(),
		User: model.User{
			ID:           101,
			Name:         "John Doe",
			Occupation:   "Developer",
			Email:        "john.doe@example.com",
			PasswordHash: "",
			Role:         "user",
		},
	}

	mockCampaignUseCase := &mocking.CampaignsUseCaseMock{}
	mockCampaignUseCase.On("CreateCampaigns", mock.AnythingOfType("model.Campaigns")).Return(mockCampaign, nil)
	controller := &campaignController{
		campaignUseCase: mockCampaignUseCase,
	}
	controller.createCampaignHandler(ginContext)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseBody dto.SingleResponse
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(suite.T(), err)
	status := responseBody.Status
	assert.Equal(suite.T(), http.StatusOK, status.Code)
	assert.Equal(suite.T(), "Campaign created successfully", status.Message)
	dataJSON, err := json.Marshal(responseBody.Data)
	assert.Nil(suite.T(), err)
	var responseCampaign model.Campaigns
	err = json.Unmarshal(dataJSON, &responseCampaign)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), responseCampaign.User.PasswordHash)
	assert.Equal(suite.T(), time.Time{}, responseCampaign.User.CreatedAt)
	assert.Equal(suite.T(), time.Time{}, responseCampaign.User.UpdatedAt)
	assert.Equal(suite.T(), mockCampaign.ID, responseCampaign.ID)
	assert.Equal(suite.T(), mockCampaign.Name, responseCampaign.Name)
	mockCampaignUseCase.AssertExpectations(suite.T())
}

// func (suite *CampaignsControllerTestSuite) TestUploadCampaignImageHandler_success() {
// 	w := httptest.NewRecorder()
// 	ginContext, _ := gin.CreateTestContext(w)
// 	campaignID := "123"
// 	ginContext.Params = append(ginContext.Params, gin.Param{Key: "campaign_id", Value: campaignID})
// 	filename := "erd.png"
// 	file, err := os.Open(" eternal-fund/images/avatars/" + filename)
// 	if err != nil {
// 		suite.T().Fatal("Failed to open test file:", err)
// 	}
// 	defer file.Close()

// 	fileHeader := &multipart.FileHeader{
// 		Filename: filename,
// 		Size:     1234,
// 	}
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, _ := writer.CreateFormFile("file", fileHeader.Filename)
// 	io.Copy(part, file)
// 	writer.Close()
// 	ginContext.Request = httptest.NewRequest(http.MethodPost, "/api/upload/"+campaignID, body)
// 	ginContext.Request.Header.Set("Content-Type", writer.FormDataContentType())
// 	input := model.CampaignImage{
// 		CampaignID:   123,
// 		FileLocation: " eternal-fund/images/campaigns/" + fileHeader.Filename,
// 	}
// 	mockCampaignImage := input
// 	mockCampaignImage.ID = 1
// 	mockCampaignUseCase := &mocking.CampaignsUseCaseMock{}
// 	mockCampaignUseCase.On("SaveCampaignImage", mock.AnythingOfType("model.CampaignImage"), mock.AnythingOfType("string")).Return(mockCampaignImage, nil)
// 	controller := &campaignController{
// 		campaignUseCase: mockCampaignUseCase,
// 	}
// 	controller.uploadCampaignImageHandler(ginContext)
// 	assert.Equal(suite.T(), http.StatusOK, w.Code)
// 	var responseBody dto.SingleResponse
// 	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
// 	assert.Nil(suite.T(), err)
// 	status := responseBody.Status
// 	assert.Equal(suite.T(), http.StatusOK, status.Code)
// 	assert.Equal(suite.T(), "Campaign image uploaded successfully", status.Message)
// 	data := responseBody.Data
// 	uploadedImage, ok := data.(model.CampaignImage)
// 	assert.True(suite.T(), ok)
// 	assert.Equal(suite.T(), mockCampaignImage.ID, uploadedImage.ID)
// 	assert.Equal(suite.T(), mockCampaignImage.CampaignID, uploadedImage.CampaignID)
// 	assert.Equal(suite.T(), mockCampaignImage.FileLocation, uploadedImage.FileLocation)
// 	mockCampaignUseCase.AssertExpectations(suite.T())
// }

func TestAuthoRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignsControllerTestSuite))
}
