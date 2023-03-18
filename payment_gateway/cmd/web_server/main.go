package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sh-valery/payment-gateway/payment_gateway/cmd/web_server/docs"
	"github.com/sh-valery/payment-gateway/payment_gateway/internal/payment"
	"github.com/sh-valery/payment-gateway/payment_gateway/internal/transport"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Payment Gateway
// @version 1.0
// @description HTTP API for payment gateway

// @contact.name   Valery
// @contact.url http://linkedin.com/in/valeryshapetin/

// @schemes   https
// @host       localhost:8080
// @BasePath  /api/v1/
func main() {
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())

	logger.Println("Connecting to database...")
	db, err := sql.Open("mysql", "username:password@tcp(db:3306)/payment_gateway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger.Println("Creating payment service...")
	repository := payment.NewPaymentRepository(db)
	service := payment.NewPaymentService(repository)
	handler := transport.NewHandler(service, logger)

	//http.HandleFunc("/api/v1/payment", handler.Payment)
	//http.HandleFunc("/api/v1/payment/status", handler.PaymentStatus)
	//
	//http.HandleFunc("/swagger/{*}", httpSwagger.Handler(
	//	httpSwagger.URL("http://localhost:8080/docs/swagger.json"), //The url pointing to API definition
	//))

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/docs/doc.json"), //The url pointing to API definition
	))
	r.Post("/api/v1/payment", handler.Payment)
	r.Get("/api/v1/payment/status", handler.PaymentStatus)

	logger.Println("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
