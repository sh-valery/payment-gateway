package payment

func maskCardNumber(cardNumber string) string {
	masked := "**** **** **** " + cardNumber[len(cardNumber)-4:]
	return masked
}
