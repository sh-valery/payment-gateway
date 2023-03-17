package main

import (
	"fmt"
	"github.com/sh-valery/payment-gateway/internal/bank"
	"github.com/sh-valery/payment-gateway/internal/transport"
	"log"
	"net/http"
)

func main() {
	handler := transport.NewHandler(bank.NewService(), nil)
	http.HandleFunc("/api/v1/deposit", handler.ChargeCard) // Update this line of code

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
