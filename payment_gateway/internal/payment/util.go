package payment

import "github.com/google/uuid"

func maskCardNumber(cardNumber string) string {
	masked := "**** **** **** " + cardNumber[len(cardNumber)-4:]
	return masked
}

type UUID struct{}

func (U UUID) New() string {
	return uuid.New().String()
}
