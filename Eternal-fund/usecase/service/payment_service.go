package service

import (
	"eternal-fund/model"
	"os"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type paymentService struct {
}

type PaymentService interface {
	GetPaymentURL(transaction model.Transaction, user model.User) (string, error)
}


func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (s *paymentService) GetPaymentURL(transaction model.Transaction, user model.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midclient.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
