package usecase

import (
	"github.com/SyamSolution/transaction-service/internal/model"
	mock "github.com/SyamSolution/transaction-service/mock"
	mock_config "github.com/SyamSolution/transaction-service/mock/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetTransactionByTransactionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock.NewMockTransactionPersister(ctrl)
	logger := mock_config.NewMockLogger(ctrl)

	uc := NewTransactionUsecase(mockTransactionRepo, logger)

	transactionID := 1
	email := "test@example.com"

	mockTransaction := model.Transaction{}
	mockDetailTransactions := []model.DetailTransaction{{}}

	mockTransactionRepo.EXPECT().GetTransactionByTransactionID(transactionID, email).Return(mockTransaction, nil)
	mockTransactionRepo.EXPECT().GetDetailTransactionByTransactionID(transactionID).Return(mockDetailTransactions, nil)

	transactionResponse, err := uc.GetTransactionByTransactionID(transactionID, email)

	assert.NotNil(t, transactionResponse)
	assert.Nil(t, err)
}

func TestGetTransactionByOrderID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock.NewMockTransactionPersister(ctrl)
	logger := mock_config.NewMockLogger(ctrl)

	uc := NewTransactionUsecase(mockTransactionRepo, logger)

	orderID := "testOrderID"

	mockTransaction := model.Transaction{}
	mockDetailTransactions := []model.DetailTransaction{{}}

	mockTransactionRepo.EXPECT().GetTransactionByOrderID(orderID).Return(mockTransaction, nil)
	mockTransactionRepo.EXPECT().GetDetailTransactionByTransactionID(mockTransaction.TransactionID).Return(mockDetailTransactions, nil)

	transactionResponse, err := uc.GetTransactionByOrderID(orderID)

	assert.NotNil(t, transactionResponse)
	assert.Nil(t, err)
}

func TestGetListTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock.NewMockTransactionPersister(ctrl)
	logger := mock_config.NewMockLogger(ctrl)

	uc := NewTransactionUsecase(mockTransactionRepo, logger)

	request := model.TransactionListRequest{}

	mockTransactions := []model.Transaction{{}}

	mockTransactionRepo.EXPECT().GetListTransaction(request).Return(mockTransactions, nil)

	transactionListResponses, err := uc.GetListTransaction(request)

	assert.NotNil(t, transactionListResponses)
	assert.Nil(t, err)
}

func TestUpdateTransactionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock.NewMockTransactionPersister(ctrl)
	logger := mock_config.NewMockLogger(ctrl)

	uc := NewTransactionUsecase(mockTransactionRepo, logger)

	orderID := "testOrderID"
	status := "testStatus"
	email := "test@example.com"

	mockTransactionRepo.EXPECT().UpdateTransactionStatus(orderID, status).Return(nil)

	err := uc.UpdateTransactionStatus(orderID, status, email)

	assert.Nil(t, err)
}
