package repository

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"log"
	"time"
)

type transactionRepo struct {
	db *sql.DB
}

func (r *transactionRepo) GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error) {
	query := "SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE campaign_id = $1"
	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.CampaignID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.Code, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepo) GetTransactionsByUserID(userID int) ([]model.Transaction, error) {
	query := "SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE user_id = $1"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.CampaignID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.Code, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepo) GetByID(id int) (model.Transaction, error) {
	var transaction model.Transaction
	query := "SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&transaction.ID, &transaction.CampaignID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.Code, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *transactionRepo) Save(transaction model.Transaction) (model.Transaction, error) {
	query := `
        INSERT INTO transactions (campaign_id, user_id, amount, status, code, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
	var id int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(query, transaction.CampaignID, transaction.UserID, transaction.Amount, transaction.Status, transaction.Code).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return transaction, err
	}
	transaction.ID = id
	transaction.CreatedAt = createdAt
	transaction.UpdatedAt = updatedAt
	return transaction, nil
}

func (r *transactionRepo) Update(transaction model.Transaction) (model.Transaction, error) {
	query := "UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	var updatedAt time.Time
	err := r.db.QueryRow(query, transaction.Status, transaction.ID).Scan(&updatedAt)
	if err != nil {
		return transaction, err
	}
	transaction.UpdatedAt = updatedAt
	return transaction, nil
}

func (r *transactionRepo) UpdatePaymentURL(transaction model.Transaction) (model.Transaction, error) {
	query := "UPDATE transactions SET payment_url = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at"
	log.Printf("Executing query: %s with payment_url: %s and id: %d", query, transaction.PaymentURL, transaction.ID)
	var updatedAt time.Time
	err := r.db.QueryRow(query, transaction.PaymentURL, transaction.ID).Scan(&updatedAt)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return transaction, err
	}
	transaction.UpdatedAt = updatedAt
	return transaction, nil
}

func (r *transactionRepo) FindAll(page int, size int) ([]model.Transaction, dto.Paging, error) {
	var listData []model.Transaction
	var rows *sql.Rows

	offset := (page - 1) * size

	rows, err := r.db.Query("SELECT * FROM transactions LIMIT $1 OFFSET $2", size, offset)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.CampaignID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.Code, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			return nil, dto.Paging{}, err
		}
		listData = append(listData, transaction)
	}

	return listData, dto.Paging{}, nil
}

func (r *transactionRepo) GetByCode(code string) (*model.Transaction, error) {
	var transaction model.Transaction
	log.Println("Querying transaction with code:", code)
	err := r.db.QueryRow("SELECT id, campaign_id, user_id, amount, status, code, payment_url, created_at, updated_at FROM transactions WHERE code = $1", code).Scan(
		&transaction.ID, &transaction.CampaignID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.Code, &transaction.PaymentURL, &transaction.CreatedAt, &transaction.UpdatedAt,
	)
	if err != nil {
		log.Println("Error querying transaction:", err)
		return nil, err
	}
	log.Println("Queried transaction:", transaction)
	return &transaction, nil
}

type TransactionRepo interface {
	GetTransactionsByCampaignID(campaignID int) ([]model.Transaction, error)
	GetTransactionsByUserID(userID int) ([]model.Transaction, error)
	GetByID(id int) (model.Transaction, error)
	Save(transaction model.Transaction) (model.Transaction, error)
	Update(transaction model.Transaction) (model.Transaction, error)
	FindAll(page int, size int) ([]model.Transaction, dto.Paging, error)
	UpdatePaymentURL(transaction model.Transaction) (model.Transaction, error)
	GetByCode(code string) (*model.Transaction, error)
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return &transactionRepo{db: db}
}
