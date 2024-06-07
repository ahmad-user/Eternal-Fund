package controller

import (
	"encoding/base64"
	"eternal-fund/middleware"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	commonresponse "eternal-fund/model/dto/common_response"
	"eternal-fund/usecase"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type TransactionController struct {
	transactionUC  usecase.TransactionUseCase
	router         *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (t *TransactionController) getCampaignTransactions(ctx *gin.Context) {
	campaignID, err := strconv.Atoi(ctx.Param("campaign_id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	transactions, err := t.transactionUC.GetTransactionsByCampaignID(campaignID)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	var data []interface{}
	for _, tx := range transactions {
		data = append(data, tx)
	}

	commonresponse.SendManyResponse(ctx, data, dto.Paging{}, "Campaign transactions retrieved successfully")
}

func (t *TransactionController) getTransactionByID(ctx *gin.Context) {
	transactionIDStr := ctx.Param("transaction_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	transaction, err := t.transactionUC.GetTransactionByID(transactionID)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	commonresponse.SendSingleResponse(ctx, transaction, "Transaction retrieved successfully")
}

func (t *TransactionController) getUserTransactions(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	transactions, err := t.transactionUC.GetTransactionsByUserID(userID)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	var data []interface{}
	for _, tx := range transactions {
		data = append(data, tx)
	}

	commonresponse.SendManyResponse(ctx, data, dto.Paging{}, "User transactions retrieved successfully")
}

func (t *TransactionController) createTransaction(ctx *gin.Context) {
	var input model.CreateTransactionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Extract userID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found"})
		return
	}

	// Set userID in input
	input.User.ID = userID.(int)

	// Debug log
	log.Printf("Handler UserID: %d", input.User.ID)

	transaction, err := t.transactionUC.CreateTransaction(input)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Memperbarui transaksi dengan URL pembayaran dari Midtrans
	// transaction.PaymentURL, err = t.transactionUC.GetPaymentURL(transaction, input.User)
	// if err != nil {
	//     commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	//     return
	// }

	commonresponse.SendSingleResponse(ctx, transaction, "Transaction created successfully")
}

func (t *TransactionController) UpdateTransaction(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("transaction_id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var input model.UpdateTransactionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	transaction, err := t.transactionUC.UpdateTransaction(transactionID, input)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	commonresponse.SendSingleResponse(ctx, transaction, "Transaction updated successfully")
}

func (t *TransactionController) getNotification(ctx *gin.Context) {
	var notificationPayload map[string]interface{}
	if err := ctx.BindJSON(&notificationPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("Notification Payload:", notificationPayload)

	transactionStatus, exists := notificationPayload["transaction_status"].(string)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload: missing transaction_status"})
		return
	}

	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload: missing order_id"})
		return
	}

	log.Println("Order ID:", orderID)
	log.Println("Transaction Status:", transactionStatus)

	// Verifikasi notifikasi dari Midtrans
	verified, err := t.verifyNotification(orderID, transactionStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify notification"})
		return
	}

	log.Println("Verified:", verified)

	if verified {
		updatedTransaction, err := t.transactionUC.UpdateTransactionStatus(orderID, transactionStatus)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction status updated", "transaction": updatedTransaction})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Notification verification failed"})
	}

}

func (t *TransactionController) verifyNotification(orderID string, transactionStatus string) (bool, error) {
	client := resty.New()
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	response, err := client.R().
		SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(serverKey+":"))).
		Get("https://api.sandbox.midtrans.com/v2/" + orderID + "/status")

	if err != nil {
		log.Println("Error verifying notification:", err)
		return false, err
	}

	log.Println("Verification response status code:", response.StatusCode())
	log.Println("Verification response body:", response.String())

	if response.StatusCode() == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func (t *TransactionController) Routing() {
	t.router.GET("/campaigns/:campaign_id/transactions", t.authMiddleware.CheckToken("user", "admin"), t.getCampaignTransactions)
	t.router.GET("/transactions/:transaction_id", t.authMiddleware.CheckToken("user", "admin"), t.getTransactionByID)
	t.router.GET("/users/:user_id/transactions", t.authMiddleware.CheckToken("user", "admin"), t.getUserTransactions)
	t.router.POST("/transactions", t.authMiddleware.CheckToken("user", "admin"), t.createTransaction)
	t.router.POST("/transactions/notification", t.getNotification)
	t.router.PUT("/transactions/:transaction_id", t.authMiddleware.CheckToken("user", "admin"), t.UpdateTransaction)
}

func NewTransactionController(transactionUc usecase.TransactionUseCase, rg *gin.RouterGroup, authMiddle middleware.AuthMiddleware) *TransactionController {
	return &TransactionController{
		transactionUC:  transactionUc,
		router:         rg,
		authMiddleware: authMiddle,
	}
}
