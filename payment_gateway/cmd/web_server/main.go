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

	// create db connection
	logger.Println("Connecting to database...")
	db, err := sql.Open("mysql", "root:pass@tcp(localhost:3306)/payment_gateway")
	//db, err := sql.Open("mysql", "username:password@tcp(db:3306)/payment_gateway")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// create service
	logger.Println("Creating payment service...")
	repository := payment.NewPaymentRepository(db)
	uuidGenerator := payment.UUID{}
	cardProcessor, err := payment.NewPartnerShipAdaptor(http.DefaultClient, "http://localhost:8081/api/v1", nil)
	if err != nil {
		log.Fatal(err)
	}
	service := payment.NewPaymentService(repository, cardProcessor, nil, uuidGenerator)
	handler := transport.NewHandler(service, logger)

	// set routing
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/docs/doc.json"), //The url pointing to API definition
	))
	r.Post("/api/v1/payment", handler.Payment)
	r.Get("/api/v1/payment/status", handler.PaymentStatus)

	// start server
	logger.Println("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
