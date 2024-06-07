package mocking

import (
    "eternal-fund/model"
    "github.com/stretchr/testify/mock"
)

type PaymentServiceMock struct {
    mock.Mock
}

func (m *PaymentServiceMock) GetPaymentURL(transaction model.Transaction, user model.User) (string, error) {
    args := m.Called(transaction, user)
    return args.String(0), args.Error(1)
}
