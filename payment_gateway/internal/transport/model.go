package transport

import "time"

type ChPaymentRequest struct {
	Source struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	} `json:"source"`
	Amount      int       `json:"amount"`
	Currency    string    `json:"currency"`
	PaymentType string    `json:"payment_type"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Capture     bool      `json:"capture"`
	CaptureOn   time.Time `json:"capture_on"`
	Customer    struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
	} `json:"customer"`
	BillingDescriptor struct {
		Name string `json:"name"`
		City string `json:"city"`
	} `json:"billing_descriptor"`
	Shipping struct {
		Address struct {
			AddressLine1 string `json:"address_line1"`
			AddressLine2 string `json:"address_line2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
			Country      string `json:"country"`
		} `json:"address"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
	} `json:"shipping"`
	Ds struct {
		Enabled    bool   `json:"enabled"`
		AttemptN3D bool   `json:"attempt_n3d"`
		Eci        string `json:"eci"`
		Cryptogram string `json:"cryptogram"`
		Xid        string `json:"xid"`
		Version    string `json:"version"`
	} `json:"3ds"`
	PreviousPaymentId string `json:"previous_payment_id"`
	Risk              struct {
		Enabled bool `json:"enabled"`
	} `json:"risk"`
	SuccessUrl string `json:"success_url"`
	FailureUrl string `json:"failure_url"`
	PaymentIp  string `json:"payment_ip"`
	Recipient  struct {
		Dob           string `json:"dob"`
		AccountNumber string `json:"account_number"`
		Zip           string `json:"zip"`
		LastName      string `json:"last_name"`
	} `json:"recipient"`
	Metadata struct {
		CouponCode string `json:"coupon_code"`
		PartnerId  int    `json:"partner_id"`
	} `json:"metadata"`
}

type ChPaymentResponse struct {
	Id              string `json:"id"`
	ActionId        string `json:"action_id"`
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
	Approved        bool   `json:"approved"`
	Status          string `json:"status"`
	AuthCode        string `json:"auth_code"`
	ResponseCode    string `json:"response_code"`
	ResponseSummary string `json:"response_summary"`
	Ds              struct {
		Downgraded bool   `json:"downgraded"`
		Enrolled   string `json:"enrolled"`
	} `json:"3ds"`
	Risk struct {
		Flagged bool `json:"flagged"`
	} `json:"risk"`
	Source struct {
		Type           string `json:"type"`
		Id             string `json:"id"`
		BillingAddress struct {
			AddressLine1 string `json:"address_line1"`
			AddressLine2 string `json:"address_line2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
			Country      string `json:"country"`
		} `json:"billing_address"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
		Last4       string `json:"last4"`
		Fingerprint string `json:"fingerprint"`
		Bin         string `json:"bin"`
	} `json:"source"`
	Customer struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
	} `json:"customer"`
	ProcessedOn time.Time `json:"processed_on"`
	Reference   string    `json:"reference"`
	Processing  struct {
		RetrievalReferenceNumber string `json:"retrieval_reference_number"`
		AcquirerTransactionId    string `json:"acquirer_transaction_id"`
		RecommendationCode       string `json:"recommendation_code"`
	} `json:"processing"`
	Eci      string `json:"eci"`
	SchemeId string `json:"scheme_id"`
	Links    struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Action struct {
			Href string `json:"href"`
		} `json:"action"`
		Void struct {
			Href string `json:"href"`
		} `json:"void"`
		Capture struct {
			Href string `json:"href"`
		} `json:"capture"`
	} `json:"_links"`
}
