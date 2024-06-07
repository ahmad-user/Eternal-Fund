package model

import (
	"time"

	"github.com/leekchan/accounting"
)

type Transaction struct {
	ID         int `json:"id"`
	CampaignID int `json:"campaign_id"`
	UserID     int `json:"user_id"`
	Amount     int `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// User       User 	
	// Campaigns  Campaigns  
}

func (t Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(t.Amount)
}

