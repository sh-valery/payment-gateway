package payment

var Config GatewayConfig

type GatewayConfig struct {
	CryptoKey []byte
	BankURL   string
	DBConnect string
}

func InitTestConfig() {
	Config = GatewayConfig{
		BankURL:   "http://bank_simulator:8081/api/v1",
		CryptoKey: []byte("mock_key_16b_len"), // rsa supports 16, 24, 32 bytes keys
		DBConnect: "root:pass@tcp(db:3306)/payment_gateway",
	}
}
