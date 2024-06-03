package helper

import (
	"bytes"
	"encoding/json"
	"github.com/SyamSolution/transaction-service/internal/model"
	"net/http"
	"os"
)

func CheckDiscount(discount model.Discount) (float32, error) {
	jsonData, err := json.Marshal(discount)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", os.Getenv("GRULE_SERVICE_URL")+"/check-discount", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var discountResponse model.ResponseDiscount
	if err := json.NewDecoder(resp.Body).Decode(&discountResponse); err != nil {
		return 0, err
	}

	return discountResponse.Data, nil
}
