package mocking

import (
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"github.com/stretchr/testify/mock"
)

type CampaignsUseCaseMock struct {
	mock.Mock
}

func (m *CampaignsUseCaseMock) CreateCampaigns(input model.Campaigns) (model.Campaigns, error) {
	args := m.Called(input)
	return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignsUseCaseMock) FindAllCampaigns(page, size int) ([]model.Campaigns, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Campaigns), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *CampaignsUseCaseMock) FindByIdCampaigns(id int) (model.Campaigns, error) {
	args := m.Called(id)
	return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignsUseCaseMock) UpdateCampaigns(id int) (model.Campaigns, error) {
	args := m.Called(id)
	return args.Get(0).(model.Campaigns), args.Error(1)
}

func (m *CampaignsUseCaseMock) DeleteCampaigns(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CampaignsUseCaseMock) SaveCampaignImage(input model.CampaignImage, fileLocation string) (model.CampaignImage, error) {
	args := m.Called(input, fileLocation)
	return args.Get(0).(model.CampaignImage), args.Error(1)
}
