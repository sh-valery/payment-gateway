package transport

import (
	"encoding/json"
	"github.com/sh-valery/payment-gateway/payment_gateway/internal/payment"
	"log"
	"net/http"
)

type Handler struct {
	svc payment.Service

	logger *log.Logger
}

func NewHandler(service payment.Service, logger *log.Logger) *Handler {
	if logger == nil {
		logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}

	return &Handler{
		svc:    service,
		logger: logger,
	}
}

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

func (h *Handler) Payment(w http.ResponseWriter, req *http.Request) {
	h.logger.Printf("Received request to charge card %+v", req)
	r := &PaymentRequest{}
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p := &payment.Payment{
		Amount:   r.Amount,
		Currency: r.Currency,
		CardInfo: payment.CardInfo{
			ExpiryMonth: r.Card.ExpiryMonth,
			ExpiryYear:  r.Card.ExpiryYear,
		},
	}

	p.CardInfo.SetCardNumber(r.Card.Number)
	p.CardInfo.SetCVV(r.Card.CVV)

	err := h.svc.ProcessPayment(p)
	if err != nil {
		h.logger.Printf("Error charging card: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := &PaymentResponse{
		ID:               p.ID,
		TrackingID:       p.TrackingID,
		MaskedCardNumber: p.CardInfo.GetMaskedNumber(),
		Status:           p.Status,
		Code:             p.StatusCode,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Printf("Error encoding response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PaymentStatus
// @Summary      Get Payment Status
// @Description  Get Payment Status by id
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200  {object}  PaymentResponse
// @Failure      400  string      string  "Bad Request"
// @Router       /payment/{id} [get]
func (h *Handler) PaymentStatus(w http.ResponseWriter, req *http.Request) {
	h.logger.Printf("PaymentStatus request", req)
	w.WriteHeader(http.StatusOK)
}
