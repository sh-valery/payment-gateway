package transport

import (
	"encoding/json"
	"github.com/sh-valery/payment-gateway/internal/payment"
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
	Payment *payment.Payment `json:"payment"`
}

type PaymentResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
}

func (h *Handler) Payment(w http.ResponseWriter, req *http.Request) {
	h.logger.Printf("Received request to charge card %+v", req)
	r := &PaymentRequest{}
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.svc.ProcessPayment(r.Payment)
	if err != nil {
		h.logger.Printf("Error charging card: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (h *Handler) PaymentStatus(w http.ResponseWriter, req *http.Request) {
	h.logger.Printf("PaymentStatus request", req)
}
