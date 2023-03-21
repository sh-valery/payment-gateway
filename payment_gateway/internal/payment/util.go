package payment

import "github.com/google/uuid"

type UUID struct{}

func (U UUID) New() string {
	return uuid.New().String()
}
