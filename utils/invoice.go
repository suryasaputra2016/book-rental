package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/suryasaputra2016/book-rental/entity"
)

func CreateInvoice(user entity.User, topupAmount float32) (*entity.Invoice, error) {
	apiKey := os.Getenv("XENDIT_APIKEY")
	apiUrl := os.Getenv("XENDIT_APIURL")

	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           topupAmount,
		"description":      "Top Up Deposit Amount of $s, $s.",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":  user.LastName + ", " + user.FirstName,
			"email": user.Email,
		},
		"currency": "IDR",
		"items": []interface{}{
			map[string]interface{}{
				"name":     "Top Up Deposit Amount",
				"quantity": 1,
				"price":    topupAmount,
			},
		},
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var resInvoice entity.Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil

}
