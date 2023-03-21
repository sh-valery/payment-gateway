package transport

import (
	"encoding/json"
	"github.com/sh-valery/payment-gateway/bank_simulator/internal/bank"
	"log"
	"net/http"
)

type Handler struct {
	svc *bank.Service

	logger *log.Logger
}

func NewHandler(service *bank.Service, logger *log.Logger) *Handler {
	if logger == nil {
		logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}

	return &Handler{
		svc:    service,
		logger: logger,
	}
}

type ChargeCardRequest struct {
	Card   bank.Card `json:"card"`
	Amount int64     `json:"amount"`
}

type ChargeCardResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
}

func (h *Handler) ChargeCard(w http.ResponseWriter, req *http.Request) {
	h.logger.Printf("Received request to charge card %+v", req)

	// Decode request body
	var reqBody ChargeCardRequest
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Printf("Received request body %+v", reqBody)

	// Validate card
	err = h.svc.ValidateCard(reqBody.Card)
	if err != nil {
		h.logger.Printf("Error validating card: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, errWrite := w.Write([]byte(err.Error()))
		if errWrite != nil {
			h.logger.Printf("Error writing response body: %v", errWrite)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Charge card
	resp, err := h.svc.ChargeCard(reqBody.Card)
	if err != nil {
		h.logger.Printf("Error charging card: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		h.logger.Printf("Error marshaling response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(respBody)
	if err != nil {
		h.logger.Printf("Error writing response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
