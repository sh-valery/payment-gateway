package payment

type Payment struct {
	ID         string
	MerchantID string
	TrackingID string // id of the payment in the payment processor
	Status     string
	StatusCode string
	Amount     int64
	Currency   string
	CardInfo   CardInfo
}

type CardInfo struct {
	ID          string
	CardToken   string
	ExpiryMonth string
	ExpiryYear  string
	// no access outside the package
	cvv        string
	cardNumber string
	holderName string
}

func (i CardInfo) GetMaskedNumber() string {
	if len(i.cardNumber) < 4 {
		return "****"
	}
	return "**** **** **** " + i.cardNumber[len(i.cardNumber)-4:]
}
