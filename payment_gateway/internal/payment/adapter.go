package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	StatusCodeSucceed = "0"
	StatusCodeFailed  = "100"
	StatusCodePending = "200"
)

type bankCard struct {
	Number     string `json:"number"`
	CVV        string `json:"cvv"`
	Expiry     string `json:"expiry"`
	HolderName string `json:"holder_name"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}

type depositCardRequest struct {
	Card   bankCard `json:"card"`
	Amount int64    `json:"amount"`
}

type bankResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Code          string `json:"code"`
}

type BankAdaptor struct {
	client *http.Client
	url    string

	logger *log.Logger
}

func NewPartnerShipAdaptor(client *http.Client, url string, logger *log.Logger) (*BankAdaptor, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}

	if logger == nil {
		logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}

	return &BankAdaptor{client: client, url: url, logger: logger}, nil
}

func (p BankAdaptor) Deposit(ctx context.Context, payment *Payment) error {
	bankRequest := generateDepositRequest(payment)

	// call bank api
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(bankRequest)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/deposit", p.url)
	res, err := p.client.Post(url, "application/json", buf)
	if err != nil {
		return fmt.Errorf("failed to call partnerships: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			p.logger.Printf("failed to close the response body of partnerships: %v", err)
		}
	}()

	// parse and validate response
	depositResponse, err := parseResponse(res)
	if err != nil {
		return err
	}

	// update payment
	payment.TrackingID = depositResponse.TransactionID
	payment.StatusCode = depositResponse.Code

	// status mapping
	switch depositResponse.Code {
	case StatusCodeSucceed:
		payment.Status = StatusSucceeded
	case StatusCodeFailed:
		payment.Status = StatusFailed
	case StatusCodePending:
		payment.Status = StatusProcessing
	default:
		return fmt.Errorf("unknown status code: %s", depositResponse.Code)
	}

	return nil
}

func parseResponse(res *http.Response) (bankResponse, error) {
	if res.StatusCode != http.StatusCreated {
		return bankResponse{}, fmt.Errorf("bad request to partnerships: %d", res.StatusCode)
	}

	var depositResponse bankResponse
	if err := json.NewDecoder(res.Body).Decode(&depositResponse); err != nil {
		return bankResponse{}, fmt.Errorf("could not decode the response body of partnerships: %w", err)
	}
	return depositResponse, nil
}

func generateDepositRequest(payment *Payment) *depositCardRequest {
	bankRequest := &depositCardRequest{
		Card: bankCard{
			Number:     payment.CardInfo.cardNumber,
			Expiry:     fmt.Sprintf("%s/%s", payment.CardInfo.ExpiryMonth, payment.CardInfo.ExpiryYear),
			CVV:        payment.CardInfo.cvv,
			HolderName: payment.CardInfo.holderName,
			Amount:     payment.Amount,
			Currency:   payment.Currency,
		},
	}
	return bankRequest
}
