package model

import (
	"time"

	"github.com/leekchan/accounting"
)

type Campaigns struct {
	ID                int             `json:"id"`
	User_id           int             `json:"user_id"`
	Name              string          `json:"name"`
	Short_description string          `json:"short_description"`
	Description       string          `json:"description"`
	Perks             string          `json:"perks"`
	Backer_count      int             `json:"backer_count"`
	Goal_amount       int             `json:"goal_amount"`
	Current_amount    int             `json:"current_amount"`
	Slug              string          `json:"slug"`
	Created_at        time.Time       `json:"created_at"`
	Updated_at        time.Time       `json:"updated_at"`
	CampaignImages    []CampaignImage `json:"campaign_images"`
	User              User       	  `json:"user"`
}

func (c Campaigns) GoalAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.Goal_amount)
}

func (c Campaigns) CurrentAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.Current_amount)
}