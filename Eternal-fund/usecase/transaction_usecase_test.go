package usecase

import (
    "eternal-fund/mocking"
    "eternal-fund/model"
    "eternal-fund/model/dto"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "testing"
)

type TransactionUseCaseTestSuite struct {
    suite.Suite
    tuc          *transactionUseCase
    transactionRepo *mocking.TransactionRepoMock
    campaignRepo *mocking.CampaignRepoMock
    paymentService *mocking.PaymentServiceMock
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
    suite.transactionRepo = new(mocking.TransactionRepoMock)
    suite.campaignRepo = new(mocking.CampaignRepoMock)
    suite.paymentService = new(mocking.PaymentServiceMock)
    suite.tuc = &transactionUseCase{
        transactionRepo: suite.transactionRepo,
        campaignRepo:    suite.campaignRepo,
        paymentService:  suite.paymentService,
    }
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByCampaignID() {
    campaignID := 1
    mockTransactions := []model.Transaction{
        {ID: 1, Amount: 1000, CampaignID: campaignID},
        {ID: 2, Amount: 2000, CampaignID: campaignID},
    }

    suite.transactionRepo.On("GetTransactionsByCampaignID", campaignID).Return(mockTransactions, nil)

    transactions, err := suite.tuc.GetTransactionsByCampaignID(campaignID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockTransactions, transactions)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByUserID() {
    userID := 1
    mockTransactions := []model.Transaction{
        {ID: 1, Amount: 1000, UserID: userID},
        {ID: 2, Amount: 2000, UserID: userID},
    }

    suite.transactionRepo.On("GetTransactionsByUserID", userID).Return(mockTransactions, nil)

    transactions, err := suite.tuc.GetTransactionsByUserID(userID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockTransactions, transactions)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionByID() {
    transactionID := 1
    mockTransaction := model.Transaction{ID: transactionID, Amount: 1000}

    suite.transactionRepo.On("GetByID", transactionID).Return(mockTransaction, nil)

    transaction, err := suite.tuc.GetTransactionByID(transactionID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockTransaction, transaction)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestCreateTransaction() {
    input := model.CreateTransactionInput{
        Amount: 1000,
        User: model.User{ID: 1},
    }
    transaction := model.Transaction{
        CampaignID: input.CampaignID,
        UserID:     input.User.ID,
        Amount:     input.Amount,
        Status:     "pending",
    }
    savedTransaction := transaction
    savedTransaction.ID = 1

    suite.transactionRepo.On("Save", transaction).Return(savedTransaction, nil)
    suite.paymentService.On("GetPaymentURL", savedTransaction, input.User).Return("http://payment.url", nil)
    suite.transactionRepo.On("UpdatePaymentURL", savedTransaction).Return(savedTransaction, nil)

    createdTransaction, err := suite.tuc.CreateTransaction(input)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), savedTransaction, createdTransaction)
    suite.transactionRepo.AssertExpectations(suite.T())
    suite.paymentService.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestUpdateTransaction() {
    transactionID := 1
    input := model.UpdateTransactionInput{Status: "paid"}
    transaction := model.Transaction{ID: transactionID, Amount: 1000, Status: "pending"}

    suite.transactionRepo.On("GetByID", transactionID).Return(transaction, nil)
    transaction.Status = input.Status
    suite.transactionRepo.On("Update", transaction).Return(transaction, nil)

    updatedTransaction, err := suite.tuc.UpdateTransaction(transactionID, input)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), transaction, updatedTransaction)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestProcessPayment() {
    input := model.TransactionNotificationInput{
        OrderID:           "1",
        PaymentType:       "credit_card",
        TransactionStatus: "capture",
        FraudStatus:       "accept",
    }
    transaction := model.Transaction{ID: 1, Amount: 1000, Status: "pending"}
    campaign := model.Campaigns{ID: 1, Backer_count: 0, Current_amount: 0}

    suite.transactionRepo.On("GetByID", 1).Return(transaction, nil)
    suite.campaignRepo.On("FindByIdCampaigns", 1).Return(campaign, nil)
    suite.transactionRepo.On("Update", transaction).Return(transaction, nil)
    suite.campaignRepo.On("UpdateCampaigns", 1).Return(campaign, nil)

    err := suite.tuc.ProcessPayment(input)
    assert.NoError(suite.T(), err)
    suite.transactionRepo.AssertExpectations(suite.T())
    suite.campaignRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestUpdateTransactionStatus() {
    orderID := "TRX-123"
    status := "paid"
    transaction := model.Transaction{ID: 1, Status: "pending", Code: orderID}

    suite.transactionRepo.On("GetByCode", orderID).Return(transaction, nil)
    transaction.Status = status
    suite.transactionRepo.On("Update", transaction).Return(transaction, nil)

    updatedTransaction, err := suite.tuc.UpdateTransactionStatus(orderID, status)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), transaction, updatedTransaction)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetAllTransactions() {
    page := 1
    size := 10
    mockTransactions := []model.Transaction{
        {ID: 1, Amount: 1000},
        {ID: 2, Amount: 2000},
    }
    paging := dto.Paging{Page: page, Size: size}

    suite.transactionRepo.On("FindAll", page, size).Return(mockTransactions, paging, nil)

    transactions, pag, err := suite.tuc.GetAllTransactions(page, size)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockTransactions, transactions)
    assert.Equal(suite.T(), paging, pag)
    suite.transactionRepo.AssertExpectations(suite.T())
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(TransactionUseCaseTestSuite))
}
