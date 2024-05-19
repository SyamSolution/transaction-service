package model

import "time"

type Transaction struct {
	TransactionID   int       `json:"transaction_id"`
	UserID          int       `json:"user_id"`
	OrderID         string    `json:"order_id"`
	TransactionDate time.Time `json:"transaction_date"`
	PaymentMethod   string    `json:"payment_method"`
	TotalAmount     float32   `json:"total_amount"`
	TotalTicket     int       `json:"total_ticket"`
	FullName        string    `json:"full_name"`
	MobileNumber    string    `json:"mobile_number"`
	Email           string    `json:"email"`
	PaymentStatus   string    `json:"payment_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DetailTransaction struct {
	DetailTransactionID int       `json:"detail_transaction_id"`
	TransactionID       int       `json:"transaction_id"`
	TicketID            int       `json:"ticket_id"`
	TicketType          string    `json:"ticket_type"`
	CountryName         string    `json:"country_name"`
	City                string    `json:"city"`
	Quantity            int       `json:"quantity"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DetailTransactionRequest struct {
	TicketID    int    `json:"ticket_id"`
	TicketType  string `json:"ticket_type"`
	CountryName string `json:"country_name"`
	City        string `json:"city"`
	Quantity    int    `json:"quantity"`
}

type TransactionRequest struct {
	PaymentMethod string                     `json:"payment_method"`
	TotalAmount   float32                    `json:"total_amount"`
	TotalTicket   int                        `json:"total_ticket"`
	DetailTicket  []DetailTransactionRequest `json:"detail_ticket"`
	PaymentStatus string                     `json:"payment_status"`
}

type DetailTransactionResponse struct {
	DetailTransactionID int    `json:"detail_transaction_id"`
	TicketID            int    `json:"ticket_id"`
	TicketType          string `json:"ticket_type"`
	CountryName         string `json:"country_name"`
	City                string `json:"city"`
	Quantity            int    `json:"quantity"`
}

type TransactionResponse struct {
	TransactionID             int                         `json:"transaction_id"`
	OrderID                   string                      `json:"order_id"`
	Email                     string                      `json:"email"`
	TransactionDate           time.Time                   `json:"transaction_date"`
	PaymentMethod             string                      `json:"payment_method"`
	TotalAmount               float32                     `json:"total_amount"`
	TotalTicket               int                         `json:"total_ticket"`
	Status                    string                      `json:"status"`
	DetailTransactionResponse []DetailTransactionResponse `json:"detail_transaction"`
}

type TransactionListRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

type TransactionListResponse struct {
	TransactionID   int       `json:"transaction_id"`
	TransactionDate time.Time `json:"transaction_date"`
	PaymentMethod   string    `json:"payment_method"`
	TotalAmount     float32   `json:"total_amount"`
	TotalTicket     int       `json:"total_ticket"`
	Status          string    `json:"status"`
}
