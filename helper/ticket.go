package helper

import (
	"encoding/json"
	"github.com/SyamSolution/transaction-service/internal/model"
	"net/http"
	"os"
	"strconv"
)

func GetTicket(continent string) ([]model.TicketResponse, error) {
	req, err := http.NewRequest("GET", os.Getenv("TICKET_MANAGEMENT_SERVICE_URL")+"/continent/tickets/"+continent, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.ResponseTicket
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func GetStockTicketGroupByContinent() ([]model.StockTicketResponse, error) {
	req, err := http.NewRequest("GET", os.Getenv("TICKET_MANAGEMENT_SERVICE_URL")+"/tickets/continent-stock", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.ResponseStockTicket
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func GetTicketEventByTicketID(ticketID int) (model.TicketEvent, error) {
	req, err := http.NewRequest("GET", os.Getenv("TICKET_MANAGEMENT_SERVICE_URL")+"/event/ticket/"+strconv.Itoa(ticketID), nil)
	if err != nil {
		return model.TicketEvent{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.TicketEvent{}, err
	}
	defer resp.Body.Close()

	var response model.ResponseTicketEvent
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return model.TicketEvent{}, err
	}

	return response.Data, nil
}
