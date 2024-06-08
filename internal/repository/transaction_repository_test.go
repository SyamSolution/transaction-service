package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyamSolution/transaction-service/internal/model"
	mock_config "github.com/SyamSolution/transaction-service/mock/config"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	transaction := model.Transaction{}
	detailTransaction := []model.DetailTransaction{{}}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO transaction").WithArgs(transaction.UserID, transaction.OrderID, transaction.TransactionDate, transaction.PaymentMethod, transaction.TotalAmount,
		transaction.TotalTicket, transaction.FullName, transaction.MobileNumber, transaction.Email, transaction.PaymentStatus, transaction.Continent,
		transaction.Discount, transaction.CreatedAt, transaction.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	for _, dt := range detailTransaction {
		mock.ExpectExec("INSERT INTO detail_transaction").WithArgs(sqlmock.AnyArg(), dt.TicketID, dt.TicketType, dt.CountryName, dt.City, dt.Quantity, dt.CreatedAt, dt.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	err = r.CreateTransaction(transaction, detailTransaction)
	if err != nil {
		t.Errorf("error was not expected while creating transaction: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionByTransactionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	transactionID := 1
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "order_id", "transaction_date", "payment_method", "total_amount", "total_ticket", "full_name",
		"mobile_number", "email", "payment_status", "continent", "created_at", "updated_at"}).
		AddRow(1, 1, "order1", time.Now(), "method1", 100, 2, "fullname1", "1234567890", "test@example.com", "status1", "continent1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT transaction_id, user_id, order_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
		mobile_number, email, payment_status, continent, created_at, updated_at FROM transaction WHERE transaction_id = \? AND email = \?`).
		WithArgs(transactionID, email).
		WillReturnRows(rows)

	_, err = r.GetTransactionByTransactionID(transactionID, email)
	if err != nil {
		t.Errorf("error was not expected while getting transaction: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionByOrderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	orderID := "order1"

	rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "order_id", "transaction_date", "payment_method", "total_amount", "total_ticket", "full_name",
		"mobile_number", "email", "payment_status", "continent", "created_at", "updated_at"}).
		AddRow(1, 1, "order1", time.Now(), "method1", 100, 2, "fullname1", "1234567890", "test@example.com", "status1", "continent1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT transaction_id, user_id, order_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
  mobile_number, email, payment_status, continent, created_at, updated_at FROM transaction WHERE order_id = \?`).
		WithArgs(orderID).
		WillReturnRows(rows)

	_, err = r.GetTransactionByOrderID(orderID)
	if err != nil {
		t.Errorf("error was not expected while getting transaction: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDetailTransactionByTransactionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	transactionID := 1

	rows := sqlmock.NewRows([]string{"detail_transaction_id", "transaction_id", "ticket_id", "ticket_type", "country_name", "city", "quantity", "created_at", "updated_at"}).
		AddRow(1, 1, 1, "type1", "country1", "city1", 1, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT detail_transaction_id, transaction_id, ticket_id, ticket_type, country_name, city, quantity, created_at, updated_at
  FROM detail_transaction WHERE transaction_id = \?`).
		WithArgs(transactionID).
		WillReturnRows(rows)

	_, err = r.GetDetailTransactionByTransactionID(transactionID)
	if err != nil {
		t.Errorf("error was not expected while getting detail transaction: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetListTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	request := model.TransactionListRequest{
		Email: "test@example.com",
	}

	rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "transaction_date", "payment_method", "total_amount", "total_ticket", "full_name",
		"mobile_number", "email", "payment_status", "continent", "created_at", "updated_at"}).
		AddRow(1, 1, time.Now(), "method1", 100, 2, "fullname1", "1234567890", "test@example.com", "status1", "continent1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT transaction_id, user_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
  mobile_number, email, payment_status, continent, created_at, updated_at FROM transaction WHERE email = \?`).
		WithArgs(request.Email).
		WillReturnRows(rows)

	_, err = r.GetListTransaction(request)
	if err != nil {
		t.Errorf("error was not expected while getting list of transactions: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	orderID := "order1"
	status := "completed"

	mock.ExpectExec(`UPDATE transaction SET payment_status = \? WHERE order_id = \?`).
		WithArgs(status, orderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = r.UpdateTransactionStatus(orderID, status)
	if err != nil {
		t.Errorf("error was not expected while updating transaction status: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDistinctContinentTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_config.NewMockLogger(ctrl)

	r := NewTransactionRepository(db, logger)

	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"continent"}).
		AddRow("continent1").
		AddRow("continent2")

	mock.ExpectQuery(`SELECT DISTINCT continent FROM transaction WHERE email = \?`).
		WithArgs(email).
		WillReturnRows(rows)

	_, err = r.GetDistinctContinentTransaction(email)
	if err != nil {
		t.Errorf("error was not expected while getting distinct continents: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
