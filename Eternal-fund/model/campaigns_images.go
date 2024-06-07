package model

import "time"

type CampaignImage struct {
	ID           int       `json:"id"`
	CampaignID   int       `json:"campaign_id"`
	FileName     string    `json:"file_name"`
	IsPrimary    int       `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FileLocation string    `form:"file_location"`
	User         User
}
