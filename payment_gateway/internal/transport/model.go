package transport

type PaymentRequest struct {
	Card struct {
		Number      string `json:"number"`
		ExpiryMonth string `json:"expiry_month"`
		ExpiryYear  string `json:"expiry_year"`
		CVV         string `json:"cvv"`
		HolderName  string `json:"holder_name"`
	} `json:"card"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type PaymentResponse struct {
	ID               string
	TrackingID       string
	MaskedCardNumber string
	Status           string `json:"status"`
	Code             string `json:"code"`
}
