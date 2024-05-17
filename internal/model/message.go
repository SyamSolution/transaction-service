package model

type Message struct {
	OrderID      string  `json:"order_id"`
	Email        string  `json:"email"`
	URL          string  `json:"url"`
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	DeadlineDate string  `json:"deadline_date"`
	Total        float32 `json:"total"`
}

type CompleteTransactionMessage struct {
	Email string `json:"email"`
}
