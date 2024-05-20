package usecase

import (
	"fmt"
	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/helper"
	"github.com/SyamSolution/transaction-service/internal/model"
	"github.com/SyamSolution/transaction-service/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

type transactionUsecase struct {
	transactionRepo repository.TransactionPersister
	logger          config.Logger
}

type TransactionExecutor interface {
	CreateTransaction(request model.TransactionRequest, user model.User) (*snap.Response, error)
	GetTransactionByTransactionID(transactionID int, email string) (model.TransactionResponse, error)
	GetTransactionByOrderID(orderID string) (model.TransactionResponse, error)
	UpdateTransactionStatus(orderID, status, email string) error
	GetListTransaction(request model.TransactionListRequest) ([]model.TransactionListResponse, error)
}

func NewTransactionUsecase(transactionRepo repository.TransactionPersister, logger config.Logger) TransactionExecutor {
	return &transactionUsecase{transactionRepo: transactionRepo, logger: logger}
}

func (uc *transactionUsecase) CreateTransaction(request model.TransactionRequest, user model.User) (*snap.Response, error) {
	orderID := fmt.Sprintf("ORDER-%s", helper.RandomOrderID(5))
	transaction := model.Transaction{
		UserID:          user.UserID,
		OrderID:         orderID,
		TransactionDate: time.Now(),
		PaymentMethod:   request.PaymentMethod,
		TotalAmount:     request.TotalAmount,
		TotalTicket:     request.TotalTicket,
		FullName:        user.FullName,
		MobileNumber:    user.PhoneNumber,
		Email:           user.Email,
		PaymentStatus:   "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	var detailTransactions []model.DetailTransaction
	for _, detail := range request.DetailTicket {
		detailTransaction := model.DetailTransaction{
			TicketID:    detail.TicketID,
			TicketType:  detail.TicketType,
			CountryName: detail.CountryName,
			City:        detail.City,
			Quantity:    detail.Quantity,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		message := model.MessageOrderTicket{
			TicketID: detailTransaction.TicketID,
			Order:    detailTransaction.Quantity,
		}

		// produce ke ticket-management-service
		if err := helper.ProduceOrderTicketMessage(message); err != nil {
			uc.logger.Error("Error when producing message order ticket", zap.Error(err))
		}

		detailTransactions = append(detailTransactions, detailTransaction)
	}

	isEligible, err := helper.CheckEligible()
	if err != nil {
		uc.logger.Error("Error when checking eligible", zap.Error(err))
		return nil, err
	}

	if !isEligible {
		return nil, fmt.Errorf("not eligible")
	}

	err = uc.transactionRepo.CreateTransaction(transaction, detailTransactions)
	if err != nil {
		uc.logger.Error("Error when creating transaction", zap.Error(err))
		return nil, err
	}

	// request ke midtrans
	var s = snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// TODO request api exchange to IDR
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(request.TotalAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.FullName,
			Email: user.Email,
			Phone: user.PhoneNumber,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	// message ke notification-service untuk send email
	// TODO tambah detail ticket
	message := model.Message{
		OrderID:      orderID,
		Email:        user.Email,
		URL:          snapResp.RedirectURL,
		Name:         user.FullName,
		Date:         time.Now().Format("02 January 2006 15:04:05"),
		DeadlineDate: time.Now().AddDate(0, 0, 1).Format("02 January 2006 15:04:05"),
		Total:        request.TotalAmount,
	}

	if err := helper.ProduceCreateTransactionMessageMail(message); err != nil {
		uc.logger.Error("Error when producing message", zap.Error(err))
	}

	return snapResp, nil
}

func (uc *transactionUsecase) GetTransactionByTransactionID(transactionID int, email string) (model.TransactionResponse, error) {
	transaction, err := uc.transactionRepo.GetTransactionByTransactionID(transactionID, email)
	if err != nil {
		uc.logger.Error("Error when getting transaction by transactionID", zap.Error(err))
		return model.TransactionResponse{}, err
	}

	detailTransactions, err := uc.transactionRepo.GetDetailTransactionByTransactionID(transactionID)
	if err != nil {
		uc.logger.Error("Error when getting detail transaction by transactionID", zap.Error(err))
		return model.TransactionResponse{}, err
	}

	var detailTransactionResponses []model.DetailTransactionResponse
	for _, detail := range detailTransactions {
		detailTransactionResponse := model.DetailTransactionResponse{
			DetailTransactionID: detail.DetailTransactionID,
			TicketID:            detail.TicketID,
			TicketType:          detail.TicketType,
			CountryName:         detail.CountryName,
			City:                detail.City,
			Quantity:            detail.Quantity,
		}
		detailTransactionResponses = append(detailTransactionResponses, detailTransactionResponse)
	}

	transactionResponse := model.TransactionResponse{
		TransactionID:             transaction.TransactionID,
		TransactionDate:           transaction.TransactionDate,
		PaymentMethod:             transaction.PaymentMethod,
		TotalAmount:               transaction.TotalAmount,
		TotalTicket:               transaction.TotalTicket,
		Status:                    transaction.PaymentStatus,
		DetailTransactionResponse: detailTransactionResponses,
	}

	return transactionResponse, nil
}

func (uc *transactionUsecase) GetTransactionByOrderID(orderID string) (model.TransactionResponse, error) {
	transaction, err := uc.transactionRepo.GetTransactionByOrderID(orderID)
	if err != nil {
		uc.logger.Error("Error when getting transaction by orderID", zap.Error(err))
		return model.TransactionResponse{}, err
	}

	detailTransactions, err := uc.transactionRepo.GetDetailTransactionByTransactionID(transaction.TransactionID)
	if err != nil {
		uc.logger.Error("Error when getting detail transaction by transactionID", zap.Error(err))
		return model.TransactionResponse{}, err
	}

	var detailTransactionResponses []model.DetailTransactionResponse
	for _, detail := range detailTransactions {
		detailTransactionResponse := model.DetailTransactionResponse{
			DetailTransactionID: detail.DetailTransactionID,
			TicketID:            detail.TicketID,
			TicketType:          detail.TicketType,
			CountryName:         detail.CountryName,
			City:                detail.City,
			Quantity:            detail.Quantity,
		}
		detailTransactionResponses = append(detailTransactionResponses, detailTransactionResponse)
	}

	transactionResponse := model.TransactionResponse{
		TransactionID:             transaction.TransactionID,
		OrderID:                   orderID,
		Email:                     transaction.Email,
		TransactionDate:           transaction.TransactionDate,
		PaymentMethod:             transaction.PaymentMethod,
		TotalAmount:               transaction.TotalAmount,
		TotalTicket:               transaction.TotalTicket,
		Status:                    transaction.PaymentStatus,
		DetailTransactionResponse: detailTransactionResponses,
	}

	return transactionResponse, nil
}

func (uc *transactionUsecase) GetListTransaction(request model.TransactionListRequest) ([]model.TransactionListResponse, error) {
	transactions, err := uc.transactionRepo.GetListTransaction(request)
	if err != nil {
		uc.logger.Error("Error when getting list transaction", zap.Error(err))
		return []model.TransactionListResponse{}, err
	}

	var transactionListResponses []model.TransactionListResponse
	for _, transaction := range transactions {
		transactionListResponse := model.TransactionListResponse{
			TransactionID:   transaction.TransactionID,
			TransactionDate: transaction.TransactionDate,
			PaymentMethod:   transaction.PaymentMethod,
			TotalAmount:     transaction.TotalAmount,
			TotalTicket:     transaction.TotalTicket,
			Status:          transaction.PaymentStatus,
		}
		transactionListResponses = append(transactionListResponses, transactionListResponse)
	}

	return transactionListResponses, nil
}

func (uc *transactionUsecase) UpdateTransactionStatus(orderID, status, email string) error {
	err := uc.transactionRepo.UpdateTransactionStatus(orderID, status)
	if err != nil {
		uc.logger.Error("Error when updating transaction status", zap.Error(err))
		return err
	}

	switch status {
	case "completed":
		log.Println("send message broker")
		message := model.CompleteTransactionMessage{Email: email}
		if err := helper.ProduceCompletedTransactionMessageMail(message); err != nil {
			uc.logger.Error("Error when producing message", zap.Error(err))
		}
	}

	return nil
}
