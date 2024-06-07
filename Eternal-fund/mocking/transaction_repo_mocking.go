package mocking

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"

	"github.com/stretchr/testify/mock"
)

type TransactionRepoMock struct {
	mock.Mock
}

func (m *TransactionRepoMock) GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error) {
	args := m.Called(campaignID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) GetTransactionsByUserID(userID int) ([]model.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) GetByID(id int) (model.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) Save(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) Update(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) UpdatePaymentURL(transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) FindAll(page int, size int) ([]model.Transaction, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]model.Transaction), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *TransactionRepoMock) GetByCode(code string) (*model.Transaction, error) {
	args := m.Called(code)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func NewTransactionRepoMock(db *sql.DB) *TransactionRepoMock {
	return &TransactionRepoMock{}
}
