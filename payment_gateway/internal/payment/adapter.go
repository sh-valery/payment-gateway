package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type bankRequest struct {
	Card struct {
		Number     string `json:"number"`
		Expiry     string `json:"expiry"`
		Cvv        string `json:"cvv"`
		HolderName string `json:"holder_name"`
	} `json:"card"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type bankResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Code          string `json:"code"`
}

type BankAdaptor struct {
	client *http.Client
	url    string
}

func NewPartnerShipAdaptor(client *http.Client, url string) (*BankAdaptor, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}

	return &BankAdaptor{client: client, url: url}, nil
}

func (p BankAdaptor) Deposit(ctx context.Context, payment *Payment) (*Payment, error) {
	bankRequest := &bankRequest{
		// todo fill in the bank request
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(bankRequest)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/deposit", p.url)
	res, err := p.client.Post(url, "application/json", buf)
	if err != nil {
		return nil, fmt.Errorf("failed to call partnerships: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request to partnerships: %d", res.StatusCode)
	}

	var pr bankResponse
	if err := json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("could not decode the response body of partnerships: %w", err)
	}

	payment.TrackingID = pr.TransactionID
	payment.Status = pr.Status

	return payment, nil
}
