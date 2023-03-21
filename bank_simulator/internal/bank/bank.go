package bank

import (
	"errors"
	"github.com/google/uuid"
)

type Card struct {
	Number     string `json:"number"`
	CVV        string `json:"cvv"`
	Expiry     string `json:"expiry"`
	HolderName string `json:"holder_name"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}

type CardChargerResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Code          string `json:"code"`
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ChargeCard(card Card) (*CardChargerResponse, error) {
	var status, code string
	switch card.Number {
	case "4242424242424242":
		status = StatusSucceeded
		code = StatusCodeSucceed
	case "4000000000000002":
		status = StatusFailed
		code = StatusCodeFailed
	case "4000000000009995": // todo: add 3ds
		status = StatusPending
		code = StatusCodePending
	}

	return &CardChargerResponse{
		TransactionID: uuid.New().String(),
		Status:        status,
		Code:          code,
	}, nil
}

func (s *Service) ValidateCard(card Card) error {
	switch {
	case card.Number == "":
		return errors.New("empty card number")
	case card.CVV == "":
		return errors.New("empty card CVV")
	case card.Expiry == "":
		return errors.New("empty card expiry")
		// some cards my have empty holder name
	}

	return nil
}
