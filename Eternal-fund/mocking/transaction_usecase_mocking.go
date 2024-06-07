package mocking

import (
	"eternal-fund/model"
	"eternal-fund/model/dto"

	"github.com/stretchr/testify/mock"
)

type TransactionUseCaseMock struct {
    mock.Mock
}

func (m *TransactionUseCaseMock) GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error) {
    args := m.Called(campaignID)
    return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) GetTransactionByID(transactionID int) (model.Transaction, error) {
    args := m.Called(transactionID)
    return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) GetTransactionsByUserID(userID int) ([]model.Transaction, error) {
    args := m.Called(userID)
    return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) CreateTransaction(input model.CreateTransactionInput) (model.Transaction, error) {
    args := m.Called(input)
    return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) UpdateTransaction(transactionID int, input model.UpdateTransactionInput) (model.Transaction, error) {
    args := m.Called(transactionID, input)
    return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) UpdateTransactionStatus(orderID string, transactionStatus string) (model.Transaction, error) {
    args := m.Called(orderID, transactionStatus)
    return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionUseCaseMock) GetAllTransactions(page, limit int) ([]model.Transaction, dto.Paging, error) {
    args := m.Called(page, limit)
    return args.Get(0).([]model.Transaction), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *TransactionUseCaseMock) GetPaymentURL(transaction model.Transaction, user model.User) (string, error) {
    args := m.Called(transaction, user)
    return args.String(0), args.Error(1)
}

func (m *TransactionUseCaseMock) ProcessPayment(input model.TransactionNotificationInput) error {
    args := m.Called(input)
    return args.Error(0)
}