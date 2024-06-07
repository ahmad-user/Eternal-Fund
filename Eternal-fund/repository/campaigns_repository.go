package repository

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"log"
	"math"
	"time"
)

type campaignsRepo struct {
	db *sql.DB
}

func (a *campaignsRepo) CreateCampaigns(campaigns model.Campaigns) (model.Campaigns, error) {
	stmt, err := a.db.Prepare(`INSERT INTO campaigns (user_id, name, short_description, description, perks,  backer_count, goal_amount,
		 current_amount, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, NOW(), NOW()) RETURNING id`)
	if err != nil {
		return model.Campaigns{}, err
	}
	defer stmt.Close()

	var campaignsID int
	err = stmt.QueryRow(campaigns.User_id, campaigns.Name, campaigns.Short_description, campaigns.Description, campaigns.Perks,
		campaigns.Backer_count, campaigns.Goal_amount, campaigns.Current_amount, campaigns.Slug).Scan(&campaignsID)
	if err != nil {
		return model.Campaigns{}, err
	}

	campaigns.ID = campaignsID

	campaigns.Created_at = time.Now()
	campaigns.Updated_at = time.Now()

	return campaigns, nil
}

func (a *campaignsRepo) FindAllCampaigns(page int, size int) ([]model.Campaigns, dto.Paging, error) {
	var listData []model.Campaigns
	var row *sql.Rows
	offset := (page - 1) * size
	var err error
	row, err = a.db.Query("SELECT  * FROM campaigns limit $1 offset $2", size, offset)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	totalRows := 0
	err = a.db.QueryRow("SELECT COUNT(*) FROM campaigns").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	for row.Next() {
		var campaigns model.Campaigns
		err := row.Scan(&campaigns.ID, &campaigns.User_id, &campaigns.Name, &campaigns.Short_description, &campaigns.Description, &campaigns.Perks, &campaigns.Backer_count,
			&campaigns.Goal_amount, &campaigns.Current_amount, &campaigns.Slug, &campaigns.Created_at, &campaigns.Updated_at)
		if err != nil {
			log.Println(err.Error())
		}
		listData = append(listData, campaigns)
	}
	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return listData, paging, nil
}

func (a *campaignsRepo) FindByIdCampaigns(id int) (model.Campaigns, error) {
	var camp model.Campaigns
	err := a.db.QueryRow("SELECT * FROM campaigns where id=$1", id).
		Scan(&camp.ID, &camp.User_id, &camp.Name, &camp.Short_description, &camp.Description, &camp.Perks, &camp.Backer_count, &camp.Goal_amount, &camp.Current_amount, &camp.Slug, &camp.Created_at, &camp.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Campaigns{}, err
		}
		return model.Campaigns{}, err
	}
	return camp, nil
}

func (a *campaignsRepo) DeleteCampaigns(id int) error {
	_, err := a.db.Exec("DELETE FROM campaigns WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (a *campaignsRepo) UpdateCampaigns(id int) (model.Campaigns, error) {
	campaign := model.Campaigns{}
	stmt, err := a.db.Prepare(`
		UPDATE campaigns 
		SET user_id = $1, name = $2, short_description = $3, description = $4,
		perks = $5, backer_count = $6, goal_amount = $7, current_amount = $8, slug = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING id, user_id, name, short_description, description, perks, backer_count, goal_amount, current_amount, slug, created_at, updated_at
	`)
	if err != nil {
		return model.Campaigns{}, err
	}
	defer stmt.Close()

	var updatedCampaign model.Campaigns
	err = stmt.QueryRow(
		campaign.User_id, campaign.Name, campaign.Short_description, campaign.Description,
		campaign.Perks, campaign.Backer_count, campaign.Goal_amount, campaign.Current_amount, campaign.Slug, campaign.ID,
	).Scan(
		&updatedCampaign.ID, &updatedCampaign.User_id, &updatedCampaign.Name, &updatedCampaign.Short_description,
		&updatedCampaign.Description, &updatedCampaign.Perks, &updatedCampaign.Backer_count, &updatedCampaign.Goal_amount,
		&updatedCampaign.Current_amount, &updatedCampaign.Slug, &updatedCampaign.Created_at, &updatedCampaign.Updated_at,
	)
	if err != nil {
		return model.Campaigns{}, err
	}

	return updatedCampaign, nil
}

func (a *campaignsRepo) FindByUserID(userID int) ([]model.Campaigns, error) {
	var campaigns []model.Campaigns
	rows, err := a.db.Query("SELECT * FROM campaigns WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var campaign model.Campaigns
		if err := rows.Scan(&campaign.ID, &campaign.User_id, &campaign.Name, &campaign.Short_description, &campaign.Description, &campaign.Perks, &campaign.Backer_count, &campaign.Goal_amount, &campaign.Current_amount, &campaign.Slug, &campaign.Created_at, &campaign.Updated_at); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}

func (a *campaignsRepo) CreateImage(campaignImage model.CampaignImage) (model.CampaignImage, error) {
	_, err := a.db.Exec("INSERT INTO campaign_images (campaign_id, file_name, is_primary) VALUES ($1, $2, $3)", campaignImage.CampaignID, campaignImage.FileName, campaignImage.IsPrimary)
	if err != nil {
		return model.CampaignImage{}, err
	}
	return campaignImage, nil
}

func (a *campaignsRepo) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	_, err := a.db.Exec("UPDATE campaign_images SET is_primary = false WHERE campaign_id = $1", campaignID)
	if err != nil {
		return false, err
	}
	return true, nil
}

type CampaignsRepo interface {
	CreateCampaigns(campaigns model.Campaigns) (model.Campaigns, error)
	FindAllCampaigns(page int, size int) ([]model.Campaigns, dto.Paging, error)
	FindByIdCampaigns(id int) (model.Campaigns, error)
	UpdateCampaigns(input int) (model.Campaigns, error)
	DeleteCampaigns(id int) error
	FindByUserID(userID int) ([]model.Campaigns, error)
	CreateImage(campaignImage model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
}

func NewCampaignsRepo(database *sql.DB) CampaignsRepo {
	return &campaignsRepo{db: database}
}
