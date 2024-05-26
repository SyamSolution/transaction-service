package model

import "time"

type TicketResponse struct {
	TicketID      int    `json:"ticket_id"`
	Type          string `json:"type"`
	Price         int    `json:"price"`
	ContinentName string `json:"continent_name"`
	Stock         int    `json:"stock"`
	CountryName   string `json:"country_name"`
	CountryCity   string `json:"country_city"`
	CountryPlace  string `json:"country_place"`
}

type StockTicketResponse struct {
	Continent string `json:"continent"`
	Stock     int    `json:"stock"`
}

type TicketEvent struct {
	TicketID     int       `json:"ticket_id"`
	Type         string    `json:"type"`
	Price        int       `json:"price"`
	Stock        int       `json:"stock"`
	Continent    string    `json:"continent"`
	CountryCity  string    `json:"country_city"`
	CountryPlace string    `json:"country_place"`
	EventName    string    `json:"event_name"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
}

type ResponseTicketEvent struct {
	Meta Meta        `json:"meta"`
	Data TicketEvent `json:"data"`
}

type ResponseTicket struct {
	Meta Meta             `json:"meta"`
	Data []TicketResponse `json:"data"`
}

type ResponseStockTicket struct {
	Meta Meta                  `json:"meta"`
	Data []StockTicketResponse `json:"data"`
}
