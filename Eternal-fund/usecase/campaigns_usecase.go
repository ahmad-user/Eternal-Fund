package usecase

import (
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"eternal-fund/repository"
	"fmt"

	"github.com/gosimple/slug"
)

type campaignsUseCase struct {
	campaignsRepo repository.CampaignsRepo
	userRepo      repository.UserRepo
}

func (a *campaignsUseCase) CreateCampaigns(input model.Campaigns) (model.Campaigns, error) {
	campaign := model.Campaigns{}
	campaign.Name = input.Name
	campaign.Short_description = input.Short_description
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.Goal_amount = input.Goal_amount
	campaign.User_id = input.User_id

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User_id)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := a.campaignsRepo.CreateCampaigns(campaign)
	if err != nil {
		return newCampaign, err
	}
	user, err := a.userRepo.FindById(input.User_id)
	if err != nil {
		return newCampaign, err
	}

	newCampaign.User = user

	return newCampaign, nil
}

func (a *campaignsUseCase) FindAllCampaigns(page int, size int) ([]model.Campaigns, dto.Paging, error) {
	campaigns, paging, err := a.campaignsRepo.FindAllCampaigns(page, size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	for i := range campaigns {
		user, err := a.userRepo.FindById(campaigns[i].User_id)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		user.PasswordHash = ""
		campaigns[i].User = user
	}

	return campaigns, paging, nil
}

func (a *campaignsUseCase) FindByIdCampaigns(inputID int) (model.Campaigns, error) {
	campaign, err := a.campaignsRepo.FindByIdCampaigns(inputID)
	if err != nil {
		return model.Campaigns{}, err
	}
	user, err := a.userRepo.FindById(campaign.User_id)
	if err != nil {
		return model.Campaigns{}, err
	}
	user.PasswordHash = ""

	campaign.User = user

	return campaign, nil

}

func (a *campaignsUseCase) UpdateCampaigns(id int) (model.Campaigns, error) {
	inputData := model.Campaigns{}
	campaign, err := a.campaignsRepo.FindByIdCampaigns(id)
	if err != nil {
		return model.Campaigns{}, err
	}

	if campaign.User_id != inputData.ID {
		return model.Campaigns{}, err
	}

	campaign.Name = inputData.Name
	campaign.Short_description = inputData.Short_description
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.Goal_amount = inputData.Goal_amount

	updatedCampaign, err := a.campaignsRepo.UpdateCampaigns(campaign.ID)
	if err != nil {
		return model.Campaigns{}, err
	}
	user, err := a.userRepo.FindById(campaign.User_id)
	if err != nil {
		return model.Campaigns{}, err
	}
	user.PasswordHash = ""

	updatedCampaign.User = user

	return updatedCampaign, nil
}

func (a *campaignsUseCase) DeleteCampaigns(id int) error {
	return a.campaignsRepo.DeleteCampaigns(id)
}

func (a *campaignsUseCase) SaveCampaignImage(input model.CampaignImage, fileLocation string) (model.CampaignImage, error) {
	inputID := model.Campaigns{
		ID: input.CampaignID,
	}
	campaign, err := a.campaignsRepo.FindByIdCampaigns(inputID.ID)
	if err != nil {
		return model.CampaignImage{}, err
	}

	if campaign.User_id != input.User.ID {
		return model.CampaignImage{}, err
	}

	isPrimary := 0
	if input.IsPrimary != 0 {
		isPrimary = 1

		_, err := a.campaignsRepo.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return model.CampaignImage{}, err
		}
	}

	campaignImage := model.CampaignImage{
		CampaignID: input.CampaignID,
		FileName:   fileLocation,
		IsPrimary:  isPrimary,
	}

	newCampaignImage, err := a.campaignsRepo.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}

type CampaignsUseCase interface {
	CreateCampaigns(input model.Campaigns) (model.Campaigns, error)
	FindAllCampaigns(page int, size int) ([]model.Campaigns, dto.Paging, error)
	FindByIdCampaigns(inputID int) (model.Campaigns, error)
	UpdateCampaigns(id int) (model.Campaigns, error)
	DeleteCampaigns(id int) error
	SaveCampaignImage(input model.CampaignImage, fileLocation string) (model.CampaignImage, error)
}

func NewCampaignsUseCase(campaignsRepo repository.CampaignsRepo, userRepo repository.UserRepo) CampaignsUseCase {
	return &campaignsUseCase{campaignsRepo: campaignsRepo, userRepo: userRepo}
}
