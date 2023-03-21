package main

import (
	"fmt"
	"github.com/sh-valery/payment-gateway/bank_simulator/internal/bank"
	"github.com/sh-valery/payment-gateway/bank_simulator/internal/transport"
	"log"
	"net/http"
)

func main() {
	handler := transport.NewHandler(bank.NewService(), nil)
	http.HandleFunc("/api/v1/deposit", handler.ChargeCard) // Update this line of code

	fmt.Printf("Starting server at port 8081\n")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
