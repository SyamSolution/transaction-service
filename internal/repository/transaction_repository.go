package repository

import (
	"database/sql"
	"github.com/SyamSolution/transaction-service/config"
	"github.com/SyamSolution/transaction-service/internal/model"
	"go.uber.org/zap"
)

type transactionRepository struct {
	DB     *sql.DB
	logger config.Logger
}

type TransactionPersister interface {
	CreateTransaction(transaction model.Transaction, detailTransaction []model.DetailTransaction) error
	GetTransactionByTransactionID(transactionID int) (model.Transaction, error)
	GetTransactionByOrderID(orderID string) (model.Transaction, error)
	GetDetailTransactionByTransactionID(transactionID int) ([]model.DetailTransaction, error)
	GetListTransaction(request model.TransactionListRequest) ([]model.Transaction, error)
	UpdateTransactionStatus(orderID string, status string) error
}

func NewTransactionRepository(DB *sql.DB, logger config.Logger) TransactionPersister {
	return &transactionRepository{DB: DB, logger: logger}
}

func (r *transactionRepository) CreateTransaction(transaction model.Transaction, detailTransaction []model.DetailTransaction) error {
	query := `INSERT INTO transaction (user_id, order_id, transaction_date, payment_method, total_amount, total_ticket, full_name, 
    		mobile_number, email, payment_status, created_at, updated_at) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`

	query2 := `INSERT INTO detail_transaction (transaction_id, ticket_type, country_name, country_code, city, quantity, created_at, updated_at)
    		VALUES (?,?,?,?,?,?,?,?)`

	tx, err := r.DB.Begin()
	if err != nil {
		r.logger.Error("Error when begin transaction", zap.Error(err))
		return err
	}

	result, err := tx.Exec(query, transaction.UserID, transaction.OrderID, transaction.TransactionDate, transaction.PaymentMethod, transaction.TotalAmount,
		transaction.TotalTicket, transaction.FullName, transaction.MobileNumber, transaction.Email, transaction.PaymentStatus, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		r.logger.Error("Error when inserting transaction", zap.Error(err))
		tx.Rollback()
		return err
	}
	idResult, _ := result.LastInsertId()

	for _, dt := range detailTransaction {
		_, err = tx.Exec(query2, idResult, dt.TicketType, dt.CountryName, dt.CountryCode, dt.City, dt.Quantity, dt.CreatedAt, dt.UpdatedAt)
		if err != nil {
			r.logger.Error("Error when inserting detail transaction", zap.Error(err))
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error when committing transaction", zap.Error(err))
		return err
	}

	return nil
}

func (r *transactionRepository) GetTransactionByTransactionID(transactionID int) (model.Transaction, error) {
	var transactions model.Transaction
	query := `SELECT transaction_id, user_id, order_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
		mobile_number, email, payment_status, created_at, updated_at FROM transaction WHERE transaction_id = ?`

	err := r.DB.QueryRow(query, transactionID).Scan(&transactions.TransactionID, &transactions.UserID, &transactions.OrderID, &transactions.TransactionDate,
		&transactions.PaymentMethod, &transactions.TotalAmount, &transactions.TotalTicket, &transactions.FullName, &transactions.MobileNumber,
		&transactions.Email, &transactions.PaymentStatus, &transactions.CreatedAt, &transactions.UpdatedAt)
	if err != nil {
		r.logger.Error("Error when scanning transaction table", zap.Error(err))
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetTransactionByOrderID(orderID string) (model.Transaction, error) {
	var transactions model.Transaction
	query := `SELECT transaction_id, user_id, order_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
		mobile_number, email, payment_status, created_at, updated_at FROM transaction WHERE order_id = ?`

	err := r.DB.QueryRow(query, orderID).Scan(&transactions.TransactionID, &transactions.UserID, &transactions.OrderID, &transactions.TransactionDate,
		&transactions.PaymentMethod, &transactions.TotalAmount, &transactions.TotalTicket, &transactions.FullName, &transactions.MobileNumber,
		&transactions.Email, &transactions.PaymentStatus, &transactions.CreatedAt, &transactions.UpdatedAt)
	if err != nil {
		r.logger.Error("Error when scanning transaction table", zap.Error(err))
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetDetailTransactionByTransactionID(transactionID int) ([]model.DetailTransaction, error) {
	var detailTransactions []model.DetailTransaction
	query := `SELECT detail_transaction_id, transaction_id, ticket_type, country_name, country_code, city, quantity, created_at, updated_at
		FROM detail_transaction WHERE transaction_id = ?`

	rows, err := r.DB.Query(query, transactionID)
	if err != nil {
		r.logger.Error("Error when querying detail transaction table", zap.Error(err))
		return detailTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var detailTransaction model.DetailTransaction
		err := rows.Scan(&detailTransaction.DetailTransactionID, &detailTransaction.TransactionID, &detailTransaction.TicketType,
			&detailTransaction.CountryName, &detailTransaction.CountryCode, &detailTransaction.City, &detailTransaction.Quantity,
			&detailTransaction.CreatedAt, &detailTransaction.UpdatedAt)
		if err != nil {
			r.logger.Error("Error when scanning detail transaction table", zap.Error(err))
			return detailTransactions, err
		}
		detailTransactions = append(detailTransactions, detailTransaction)
	}

	return detailTransactions, nil
}

func (r *transactionRepository) GetListTransaction(request model.TransactionListRequest) ([]model.Transaction, error) {
	var transactions []model.Transaction
	query := `SELECT transaction_id, user_id, transaction_date, payment_method, total_amount, total_ticket, full_name,
		mobile_number, email, payment_status, created_at, updated_at FROM transaction 
		WHERE email = ?`

	if request.Status != "" {
		query += " AND payment_status = '" + request.Status + "'"
	}

	//if request.StartDate.IsZero() && !request.EndDate.IsZero() {
	//	startDate := request.StartDate.Format("2006-01-02")
	//	endDate := request.EndDate.Format("2006-01-02")
	//	query += " AND transaction_date BETWEEN '" + startDate + "' AND '" + endDate + "'"
	//}

	rows, err := r.DB.Query(query, request.Email)
	if err != nil {
		r.logger.Error("Error when querying transaction table", zap.Error(err))
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.TransactionID, &transaction.UserID, &transaction.TransactionDate,
			&transaction.PaymentMethod, &transaction.TotalAmount, &transaction.TotalTicket, &transaction.FullName, &transaction.MobileNumber,
			&transaction.Email, &transaction.PaymentStatus, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			r.logger.Error("Error when scanning transaction table", zap.Error(err))
			return transactions, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) UpdateTransactionStatus(orderID string, status string) error {
	query := `UPDATE transaction SET payment_status = ? WHERE order_id = ?`

	_, err := r.DB.Exec(query, status, orderID)
	if err != nil {
		r.logger.Error("Error when updating transaction status", zap.Error(err))
		return err
	}

	return nil
}
