package controller

import (
    "encoding/json"
    "eternal-fund/mocking"
    "eternal-fund/model"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type TransactionControllerTestSuite struct {
    suite.Suite
    router *gin.Engine
    tuc    *mocking.TransactionUseCaseMock
    amm    *mocking.AuthMiddlewareMock
}

func (suite *TransactionControllerTestSuite) SetupTest() {
    suite.tuc = new(mocking.TransactionUseCaseMock)
    suite.amm = new(mocking.AuthMiddlewareMock)
    suite.router = gin.Default()
    gin.SetMode(gin.TestMode)
    rg := suite.router.Group("/api/v1")

    transactionController := NewTransactionController(suite.tuc, rg, suite.amm)
    transactionController.Routing()
}

func (suite *TransactionControllerTestSuite) TestGetCampaignTransactions() {
    campaignID := 1
    mockTransactions := []model.Transaction{
        {ID: 1, Amount: 1000, CampaignID: campaignID},
        {ID: 2, Amount: 2000, CampaignID: campaignID},
    }

    suite.tuc.On("GetTransactionsByCampaignID", campaignID).Return(mockTransactions, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/campaigns/1/transactions", nil)
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestGetTransactionByID() {
    transactionID := 1
    mockTransaction := model.Transaction{ID: transactionID, Amount: 1000}

    suite.tuc.On("GetTransactionByID", transactionID).Return(mockTransaction, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/transactions/1", nil)
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestGetUserTransactions() {
    userID := 1
    mockTransactions := []model.Transaction{
        {ID: 1, Amount: 1000, UserID: userID},
        {ID: 2, Amount: 2000, UserID: userID},
    }

    suite.tuc.On("GetTransactionsByUserID", userID).Return(mockTransactions, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/users/1/transactions", nil)
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction() {
    input := model.CreateTransactionInput{
        Amount: 1000,
        User: model.User{ID: 1},
    }
    transaction := model.Transaction{ID: 1, Amount: 1000}

    suite.tuc.On("CreateTransaction", input).Return(transaction, nil)

    body, _ := json.Marshal(input)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/transactions", strings.NewReader(string(body)))
    req.Header.Set("Content-Type", "application/json")
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestUpdateTransaction() {
    transactionID := 1
    input := model.UpdateTransactionInput{Amount: 2000}
    transaction := model.Transaction{ID: transactionID, Amount: 2000}

    suite.tuc.On("UpdateTransaction", transactionID, input).Return(transaction, nil)

    body, _ := json.Marshal(input)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/api/v1/transactions/1", strings.NewReader(string(body)))
    req.Header.Set("Content-Type", "application/json")
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestGetNotification() {
    notificationPayload := map[string]interface{}{
        "transaction_status": "settlement",
        "order_id":           "order-123",
    }
    body, _ := json.Marshal(notificationPayload)

    suite.tuc.On("UpdateTransactionStatus", "order-123", "settlement").Return(model.Transaction{ID: 1, Status: "settlement"}, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/transactions/notification", strings.NewReader(string(body)))
    req.Header.Set("Content-Type", "application/json")
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
    suite.tuc.AssertExpectations(suite.T())
}

func TestTransactionControllerTestSuite(t *testing.T) {
    suite.Run(t, new(TransactionControllerTestSuite))
}
