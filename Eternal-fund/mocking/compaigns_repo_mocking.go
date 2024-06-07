package mocking

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"

	"github.com/stretchr/testify/mock"
)

type CampaignRepoMock struct {
    mock.Mock
}

func (m *CampaignRepoMock) FindByIdCampaigns(id int) (model.Campaigns, error) {
    args := m.Called(id)
    return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignRepoMock) UpdateCampaigns(id int) (model.Campaigns, error) {
    args := m.Called(id)
    return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignRepoMock) FindAllCampaigns(page int, size int) ([]model.Campaigns, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Campaigns), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *CampaignRepoMock) CreateCampaigns(campaigns model.Campaigns) (model.Campaigns, error) {
	args := m.Called(campaigns)
	return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignRepoMock) CreateImage(image model.CampaignImage) (model.CampaignImage, error) {
    args := m.Called(image)
    return args.Get(0).(model.CampaignImage), args.Error(1)
}

func (m *CampaignRepoMock) DeleteCampaigns(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *CampaignRepoMock) FindByUserID(userID int) ([]model.Campaigns, error) {
    args := m.Called(userID)
    return args.Get(0).([]model.Campaigns), args.Error(1)
}

func (m *CampaignRepoMock) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	args := m.Called(campaignID)
	return args.Bool(0), args.Error(1)
}



func NewCampaignRepoMock(db *sql.DB) *CampaignRepoMock {
	return &CampaignRepoMock{}
}

