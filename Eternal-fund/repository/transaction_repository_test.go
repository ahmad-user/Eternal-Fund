package repository

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var expectedTransaction = model.Transaction{
	ID:         23,
	CampaignID: 3,
	UserID:     1,
	Amount:     100000000,
	Status:     "settlement",
	Code:       "TRX-1717468985",
	PaymentURL: "https://app.sandbox.midtrans.com/snap/v4/redirection/796f4f19-e122-4b13-b537-6911a38a1b37",
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
}

type TransactionRepoTestSuite struct {
	suite.Suite
	mockDB          *sql.DB
	mockSql         sqlmock.Sqlmock
	transactionRepo TransactionRepo
}

func (suite *TransactionRepoTestSuite) SetupTest() {
	mockDB, mockSql, _ := sqlmock.New()
	suite.mockDB = mockDB
	suite.mockSql = mockSql
	suite.transactionRepo = NewTransactionRepo(suite.mockDB)
}

func (suite *TransactionRepoTestSuite) TestGetTransactionsByCampaignID_Success() {
	campaignID := 3

	rows := sqlmock.NewRows([]string{"id", "campaign_id", "user_id", "amount", "status", "code", "payment_url", "created_at", "updated_at"}).
		AddRow(expectedTransaction.ID, expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code, expectedTransaction.PaymentURL, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE campaign_id = $1`)).
		WithArgs(campaignID).
		WillReturnRows(rows)

	actualTransactions, err := suite.transactionRepo.GetTransactionsByCampaignID(campaignID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []model.Transaction{expectedTransaction}, actualTransactions)
}

func (suite *TransactionRepoTestSuite) TestGetTransactionsByCampaignID_Fail() {
	campaignID := 3

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE campaign_id = $1`)).
		WithArgs(campaignID).
		WillReturnError(fmt.Errorf("error"))

	actualTransactions, err := suite.transactionRepo.GetTransactionsByCampaignID(campaignID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), actualTransactions)
}

func (suite *TransactionRepoTestSuite) TestGetTransactionsByUserID_Success() {
	userID := 1

	rows := sqlmock.NewRows([]string{"id", "campaign_id", "user_id", "amount", "status", "code", "payment_url", "created_at", "updated_at"}).
		AddRow(expectedTransaction.ID, expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code, expectedTransaction.PaymentURL, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	actualTransactions, err := suite.transactionRepo.GetTransactionsByUserID(userID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []model.Transaction{expectedTransaction}, actualTransactions)
}

func (suite *TransactionRepoTestSuite) TestGetTransactionsByUserID_Fail() {
	userID := 1

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error"))

	actualTransactions, err := suite.transactionRepo.GetTransactionsByUserID(userID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), actualTransactions)
}

func (suite *TransactionRepoTestSuite) TestGetByID_Success() {
	rows := sqlmock.NewRows([]string{"id", "campaign_id", "user_id", "amount", "status", "code", "payment_url", "created_at", "updated_at"}).
		AddRow(expectedTransaction.ID, expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code, expectedTransaction.PaymentURL, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE id = $1`)).
		WithArgs(expectedTransaction.ID).
		WillReturnRows(rows)

	actualTransaction, actualError := suite.transactionRepo.GetByID(expectedTransaction.ID)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedTransaction, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestGetByID_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE id = $1`)).
		WithArgs(expectedTransaction.ID).
		WillReturnError(fmt.Errorf("error"))

	actualTransaction, actualError := suite.transactionRepo.GetByID(expectedTransaction.ID)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Transaction{}, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestFindAll_Success() {
	page := 1
	size := 2
	offset := (page - 1) * size

	expectedTransactions := []model.Transaction{
		{ID: 1, CampaignID: 1, UserID: 1, Amount: 10000000, Status: "pending", Code: "TRX-1", PaymentURL: "https://payment-url.com/1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, CampaignID: 2, UserID: 2, Amount: 20000000, Status: "success", Code: "TRX-2", PaymentURL: "https://payment-url.com/2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "campaign_id", "user_id", "amount", "status", "code", "payment_url", "created_at", "updated_at"}).
		AddRow(expectedTransactions[0].ID, expectedTransactions[0].CampaignID, expectedTransactions[0].UserID, expectedTransactions[0].Amount,
			expectedTransactions[0].Status, expectedTransactions[0].Code, expectedTransactions[0].PaymentURL, expectedTransactions[0].CreatedAt, expectedTransactions[0].UpdatedAt).
		AddRow(expectedTransactions[1].ID, expectedTransactions[1].CampaignID, expectedTransactions[1].UserID, expectedTransactions[1].Amount,
			expectedTransactions[1].Status, expectedTransactions[1].Code, expectedTransactions[1].PaymentURL, expectedTransactions[1].CreatedAt, expectedTransactions[1].UpdatedAt)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM transactions LIMIT $1 OFFSET $2`)).
		WithArgs(size, offset).
		WillReturnRows(rows)

	totalRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM transactions")).
		WillReturnRows(totalRows)

	actualTransactions, actualPaging, actualError := suite.transactionRepo.FindAll(page, size)

	assert.NoError(suite.T(), actualError)
	assert.Equal(suite.T(), dto.Paging{Page: page, Size: size, TotalRows: 5, TotalPages: 3}, actualPaging)
	assert.Equal(suite.T(), expectedTransactions, actualTransactions)
}

func (suite *TransactionRepoTestSuite) TestFindAll_Fail() {
	size := 2
	page := 1
	offset := (page - 1) * size

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM transactions LIMIT $1 OFFSET $2`)).
		WithArgs(size, offset).
		WillReturnError(fmt.Errorf("error fetching transactions"))

	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM transactions")).
		WillReturnError(fmt.Errorf("error counting total rows"))

	actualTransactions, actualPaging, actualError := suite.transactionRepo.FindAll(page, size)

	assert.Error(suite.T(), actualError)
	assert.Empty(suite.T(), actualTransactions)
	assert.Equal(suite.T(), dto.Paging{}, actualPaging)
}

func (suite *TransactionRepoTestSuite) TestGetByCode_Success() {
	code := "TRX-1717468985"

	rows := sqlmock.NewRows([]string{"id", "campaign_id", "user_id", "amount", "status", "code", "payment_url", "created_at", "updated_at"}).
		AddRow(expectedTransaction.ID, expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code, expectedTransaction.PaymentURL, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE code = $1`)).
		WithArgs(code).
		WillReturnRows(rows)

	actualTransaction, actualError := suite.transactionRepo.GetByCode(code)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), &expectedTransaction, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestGetByCode_Fail() {
	code := "TRX-1717468985"

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE code = $1`)).
		WithArgs(code).
		WillReturnError(fmt.Errorf("error"))

	actualTransaction, actualError := suite.transactionRepo.GetByCode(code)

	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestSave_Success() {
	expectedQuery := `INSERT INTO transactions \(campaign_id, user_id, amount, status, code, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\) RETURNING id, created_at, updated_at`

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(expectedTransaction.ID, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt))

	actualTransaction, err := suite.transactionRepo.Save(expectedTransaction)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTransaction, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestSave_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`INSERT INTO transactions (campaign_id, user_id, amount, status, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`)).
		WithArgs(expectedTransaction.CampaignID, expectedTransaction.UserID, expectedTransaction.Amount,
			expectedTransaction.Status, expectedTransaction.Code).
		WillReturnError(fmt.Errorf("error"))
	actualTransaction, err := suite.transactionRepo.Save(expectedTransaction)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Transaction{}, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestUpdate_Success() {
	expectedQuery := "UPDATE transactions SET status = $1, updated_at = NOW\\(\\) WHERE id = $2 RETURNING updated_at"

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(expectedTransaction.Status, expectedTransaction.ID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).
			AddRow(expectedTransaction.UpdatedAt))

	actualTransaction, err := suite.transactionRepo.Update(expectedTransaction)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTransaction, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestUpdate_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`)).
		WithArgs(expectedTransaction.Status, expectedTransaction.ID).
		WillReturnError(fmt.Errorf("error"))

	actualTransaction, err := suite.transactionRepo.Update(expectedTransaction)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Transaction{}, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestUpdatePaymentURL_Success() {
	expectedQuery := "UPDATE transactions SET payment_url = $1, updated_at = NOW\\(\\) WHERE id = $2 RETURNING updated_at"

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(expectedTransaction.PaymentURL, expectedTransaction.ID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).
			AddRow(expectedTransaction.UpdatedAt))

	actualTransaction, err := suite.transactionRepo.UpdatePaymentURL(expectedTransaction)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTransaction, actualTransaction)
}

func (suite *TransactionRepoTestSuite) TestUpdatePaymentURL_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`UPDATE transactions SET payment_url = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`)).
		WithArgs(expectedTransaction.PaymentURL, expectedTransaction.ID).
		WillReturnError(fmt.Errorf("error"))

	actualTransaction, err := suite.transactionRepo.UpdatePaymentURL(expectedTransaction)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Transaction{}, actualTransaction)
}

func TestTransactionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepoTestSuite))
}
