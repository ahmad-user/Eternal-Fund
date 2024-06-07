package usecase

import (
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"eternal-fund/repository"
	"eternal-fund/usecase/service"
	"fmt"
	"log"
	"strconv"
	"time"
)

type transactionUseCase struct {
	transactionRepo repository.TransactionRepo
	campaignRepo    repository.CampaignsRepo
	paymentService  service.PaymentService
}

func (uc *transactionUseCase) GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error) {
	return uc.transactionRepo.GetTransactionsByCampaignID(campaignID)
}
func (uc *transactionUseCase) GetTransactionsByUserID(userID int) ([]model.Transaction, error) {
	return uc.transactionRepo.GetTransactionsByUserID(userID)
}

func (u *transactionUseCase) GetTransactionByID(id int) (model.Transaction, error) {
	return u.transactionRepo.GetByID(id)
}

func (uc *transactionUseCase) GetPaymentURL(transaction model.Transaction, user model.User) (string, error) {
	return uc.paymentService.GetPaymentURL(transaction, user)
}

func (uc *transactionUseCase) CreateTransaction(input model.CreateTransactionInput) (model.Transaction, error) {
	transaction := model.Transaction{
		CampaignID: input.CampaignID,
		UserID:     input.User.ID,
		Amount:     input.Amount,
		Status:     "pending",
		Code:       "TRX-" + strconv.Itoa(int(time.Now().Unix())),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Log input data
	fmt.Printf("Creating transaction: %+v\n", transaction)

	savedTransaction, err := uc.transactionRepo.Save(transaction)
	if err != nil {
		fmt.Printf("Error saving transaction: %v\n", err)
		return model.Transaction{}, err
	}
	fmt.Printf("Transaction saved: %+v\n", savedTransaction)

	paymentURL, err := uc.paymentService.GetPaymentURL(savedTransaction, input.User)
	if err != nil {
		fmt.Printf("Error getting payment URL: %v\n", err)
		return model.Transaction{}, err
	}

	savedTransaction.PaymentURL = paymentURL
	updatedTransaction, err := uc.transactionRepo.UpdatePaymentURL(savedTransaction)
	if err != nil {
		fmt.Printf("Error updating transaction with payment URL: %v\n", err)
		return model.Transaction{}, err
	}
	fmt.Printf("Transaction updated with payment URL: %+v\n", updatedTransaction)

	return updatedTransaction, nil
}

func (u *transactionUseCase) ProcessPayment(input model.TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := u.transactionRepo.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := u.transactionRepo.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := u.campaignRepo.FindByIdCampaigns(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.Backer_count = campaign.Backer_count + 1
		campaign.Current_amount = campaign.Current_amount + updatedTransaction.Amount

		_, err := u.campaignRepo.UpdateCampaigns(campaign.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *transactionUseCase) UpdateTransaction(transactionID int, input model.UpdateTransactionInput) (model.Transaction, error) {
	// Fetch the existing transaction
	transaction, err := uc.transactionRepo.GetByID(transactionID)
	if err != nil {
		return model.Transaction{}, err
	}

	// Update the transaction fields
	// transaction.Amount = input.Amount
	transaction.Status = input.Status

	// Save the updated transaction
	updatedTransaction, err := uc.transactionRepo.Update(transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	return updatedTransaction, nil
}

func (u *transactionUseCase) GetAllTransactions(page int, size int) ([]model.Transaction, dto.Paging, error) {
	return u.transactionRepo.FindAll(page, size)
}

func (uc *transactionUseCase) UpdateTransactionStatus(orderID, status string) (model.Transaction, error) {
	log.Println("Fetching transaction with order ID:", orderID)
	transaction, err := uc.transactionRepo.GetByCode(orderID)
	if err != nil {
		log.Println("Error fetching transaction by code:", err)
		return model.Transaction{}, err
	}

	log.Println("Fetched transaction:", transaction)

	transaction.Status = status
	log.Println("Updating transaction status to:", status)
	updatedTransaction, err := uc.transactionRepo.Update(*transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	return updatedTransaction, nil
}

type TransactionUseCase interface {
	GetPaymentURL(transaction model.Transaction, user model.User) (string, error)
	GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error)
	GetTransactionsByUserID(userID int) ([]model.Transaction, error)
	GetTransactionByID(id int) (model.Transaction, error)
	CreateTransaction(input model.CreateTransactionInput) (model.Transaction, error)
	UpdateTransaction(transactionID int, input model.UpdateTransactionInput) (model.Transaction, error)
	UpdateTransactionStatus(orderID, status string) (model.Transaction, error)
	ProcessPayment(input model.TransactionNotificationInput) error
	GetAllTransactions(page int, size int) ([]model.Transaction, dto.Paging, error)
}

type PaymentService interface {
	GetPaymentURL(transaction model.Transaction, user model.User) (string, error)
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepo, campaignRepo repository.CampaignsRepo, paymentService service.PaymentService) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		campaignRepo:    campaignRepo,
		paymentService:  paymentService,
	}
}
