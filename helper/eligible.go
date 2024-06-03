package helper

import (
	"bytes"
	"encoding/json"
	"github.com/SyamSolution/transaction-service/internal/model"
	"net/http"
	"os"
	"time"
)

func CheckEligible() (bool, error) {
	eligible := model.Eligible{
		Day:   time.Now().Format("Monday"),
		Month: time.Now().Format("January"),
		Year:  time.Now().Format("2006"),
	}
	jsonData, err := json.Marshal(eligible)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", os.Getenv("GRULE_SERVICE_URL")+"/check-eligible", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var eligibleResponse model.ResponseEligible
	if err := json.NewDecoder(resp.Body).Decode(&eligibleResponse); err != nil {
		return false, err
	}

	return eligibleResponse.Data, nil
}
