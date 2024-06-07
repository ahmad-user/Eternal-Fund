package usecase

import (
	"eternal-fund/mocking"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CampaignUseCaseTestSuite struct {
	suite.Suite
	cuc          *campaignsUseCase
	campaignRepo *mocking.CampaignRepoMock
	userRepo     *mocking.UserRepoMock
}

func (suite *CampaignUseCaseTestSuite) SetupTest() {
	suite.campaignRepo = new(mocking.CampaignRepoMock)
	suite.userRepo = new(mocking.UserRepoMock)
	suite.cuc = &campaignsUseCase{
		campaignsRepo: suite.campaignRepo,
		userRepo:      suite.userRepo,
	}
}

func (suite *CampaignUseCaseTestSuite) TestCreateCampaigns() {
	input := model.Campaigns{
		Name:              "Test Campaign",
		Short_description: "Short Description",
		Description:       "Long Description",
		Goal_amount:       100000,
		User_id:           1,
	}
	savedCampaign := input
	savedCampaign.ID = 1

	suite.campaignRepo.On("CreateCampaigns", input).Return(savedCampaign, nil)
	suite.userRepo.On("FindById", input.User_id).Return(model.User{}, nil)

	createdCampaign, err := suite.cuc.CreateCampaigns(input)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), savedCampaign, createdCampaign)
	suite.campaignRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CampaignUseCaseTestSuite) TestFindAllCampaigns() {
	page := 1
	size := 10
	mockCampaigns := []model.Campaigns{
		{ID: 1, Name: "Campaign 1"},
		{ID: 2, Name: "Campaign 2"},
	}
	mockPaging := dto.Paging{
		Page: page,
		Size: size,
	}

	suite.campaignRepo.On("FindAllCampaigns", page, size).Return(mockCampaigns, mockPaging, nil)
	suite.userRepo.On("FindById", mockCampaigns[0].User_id).Return(model.User{}, nil)
	suite.userRepo.On("FindById", mockCampaigns[1].User_id).Return(model.User{}, nil)

	campaigns, paging, err := suite.cuc.FindAllCampaigns(page, size)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCampaigns, campaigns)
	assert.Equal(suite.T(), mockPaging, paging)
	suite.campaignRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CampaignUseCaseTestSuite) TestFindByIdCampaigns() {
	campaignID := 1
	mockCampaign := model.Campaigns{ID: campaignID, Name: "Test Campaign"}

	suite.campaignRepo.On("FindByIdCampaigns", campaignID).Return(mockCampaign, nil)
	suite.userRepo.On("FindById", mockCampaign.User_id).Return(model.User{}, nil)

	campaign, err := suite.cuc.FindByIdCampaigns(campaignID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCampaign, campaign)
	suite.campaignRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CampaignUseCaseTestSuite) TestUpdateCampaigns() {
	campaignID := 1
	mockCampaign := model.Campaigns{ID: campaignID, Name: "Test Campaign"}

	suite.campaignRepo.On("FindByIdCampaigns", campaignID).Return(mockCampaign, nil)
	suite.userRepo.On("FindById", mockCampaign.User_id).Return(model.User{}, nil)
	suite.campaignRepo.On("UpdateCampaigns", campaignID).Return(mockCampaign, nil)

	updatedCampaign, err := suite.cuc.UpdateCampaigns(campaignID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockCampaign, updatedCampaign)
	suite.campaignRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *CampaignUseCaseTestSuite) TestDeleteCampaigns() {
	campaignID := 1

	suite.campaignRepo.On("DeleteCampaigns", campaignID).Return(nil)

	err := suite.cuc.DeleteCampaigns(campaignID)
	assert.NoError(suite.T(), err)
	suite.campaignRepo.AssertExpectations(suite.T())
}

func (suite *CampaignUseCaseTestSuite) TestSaveCampaignImage() {
	input := model.CampaignImage{
		CampaignID: 1,
		User:       model.User{ID: 1},
		IsPrimary:  1,
	}
	fileLocation := "/path/to/image.jpg"

	suite.campaignRepo.On("FindByIdCampaigns", input.CampaignID).Return(model.Campaigns{}, nil)
	suite.campaignRepo.On("MarkAllImagesAsNonPrimary", input.CampaignID).Return(nil)
	suite.campaignRepo.On("CreateImage", input).Return(input, nil)

	newCampaignImage, err := suite.cuc.SaveCampaignImage(input, fileLocation)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), input, newCampaignImage)
	suite.campaignRepo.AssertExpectations(suite.T())
}

func TestCampaignUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignUseCaseTestSuite))
}
