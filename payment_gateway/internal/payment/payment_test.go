package payment_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/sh-valery/payment-gateway/payment_gateway/internal/payment"
	"github.com/sh-valery/payment-gateway/payment_gateway/internal/payment/mock"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestServiceImpl_ProcessPayment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	succeedCard := &payment.CardInfo{}
	succeedCard.SetCardNumber("4242424242424242")

	failedCard := &payment.CardInfo{}
	failedCard.SetCardNumber("4000000000000002")

	processingCard := &payment.CardInfo{}
	processingCard.SetCardNumber("4000000000009995")

	tests := []struct {
		name    string
		service payment.Service
		args    *payment.Payment
		want    *payment.Payment
		wantErr bool
	}{
		{
			name:    "succeed payment",
			service: generateMockService(ctrl),
			args: &payment.Payment{
				MerchantID: "testMerchantID",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   succeedCard,
			},
			want: &payment.Payment{
				ID:         "mockID",
				MerchantID: "testMerchantID",
				TrackingID: "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
				Status:     payment.StatusSucceeded,
				StatusCode: "0",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   succeedCard,
			},
		},
		{
			name:    "failed payment",
			service: generateMockService(ctrl),
			args: &payment.Payment{
				MerchantID: "testMerchantID",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   failedCard,
			},
			want: &payment.Payment{
				ID:         "mockID",
				MerchantID: "testMerchantID",
				TrackingID: "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
				Status:     payment.StatusFailed,
				StatusCode: "100",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   failedCard,
			},
		},
		{
			name:    "processing payment",
			service: generateMockService(ctrl),
			args: &payment.Payment{
				MerchantID: "testMerchantID",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   processingCard,
			},
			want: &payment.Payment{
				ID:         "mockID",
				MerchantID: "testMerchantID",
				TrackingID: "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
				Status:     payment.StatusProcessing,
				StatusCode: "200",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   processingCard,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.service.ProcessPayment(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaymentDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args, tt.want) {
				t.Errorf("GetPaymentDetails() got = %v, want %v", tt.args, tt.want)
			}
		})
	}
}

func TestServiceImpl_GetPaymentDetails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	succeedCard := &payment.CardInfo{}
	succeedCard.SetCardNumber("4242424242424242")

	failedCard := &payment.CardInfo{}
	failedCard.SetCardNumber("4000000000000002")

	processingCard := &payment.CardInfo{}
	processingCard.SetCardNumber("4000000000009995")

	tests := []struct {
		name    string
		service payment.Service
		args    struct {
			paymentID  string
			merchantID string
		}
		want    *payment.Payment
		wantErr bool
	}{
		{
			name:    "succeed payment",
			service: generateMockService(ctrl),
			args: struct {
				paymentID  string
				merchantID string
			}{
				paymentID:  "mockID",
				merchantID: "testMerchantID",
			},
			want: &payment.Payment{
				ID:         "mockID",
				MerchantID: "testMerchantID",
				TrackingID: "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
				Status:     payment.StatusSucceeded,
				StatusCode: "0",
				Amount:     99,
				Currency:   "USD",
				CardInfo:   succeedCard,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := tt.service.GetPaymentDetails(tt.args.paymentID, tt.args.merchantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaymentDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("GetPaymentDetails() got = %v, want %v", tt.args, tt.want)
			}
		})
	}
}

func generateMockService(ctrl *gomock.Controller) payment.Service {
	// Mock repository behaviour
	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().Create(gomock.Any()).Return(nil).AnyTimes()
	mockRepository.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockRepository.EXPECT().GetByID("mockID").Return(&payment.Payment{
		ID:         "mockID",
		MerchantID: "testMerchantID",
		TrackingID: "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
		Status:     payment.StatusSucceeded,
		StatusCode: "0",
		Amount:     99,
		Currency:   "USD",
		CardInfo: &payment.CardInfo{
			ID: "mockCardID",
		},
	}, nil).AnyTimes()

	mockCardRepository := mock.NewMockCardRepository(ctrl)
	mockCardRepository.EXPECT().SaveCardInfo(gomock.Any()).Return(nil).AnyTimes()
	card := &payment.CardInfo{}
	card.SetCardNumber("4242424242424242")
	mockCardRepository.EXPECT().GetCardByID("mockCardID").Return(card, nil).AnyTimes()

	cardProcessor, err := payment.NewPartnerShipAdaptor(http.DefaultClient, "http://bank_simulator:8081/api/v1", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Mock bank behaviour
	httpmock.RegisterMatcherResponder("POST", "http://bank_simulator:8081/api/v1/deposit",
		httpmock.BodyContainsString(`"number":"4242424242424242"`),
		httpmock.NewStringResponder(http.StatusCreated, `{
								  "transaction_id": "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
								  "status": "succeeded",
								  "code": "0"
								}`))
	httpmock.RegisterMatcherResponder("POST", "http://bank_simulator:8081/api/v1/deposit",
		httpmock.BodyContainsString(`"number":"4000000000000002"`),
		httpmock.NewStringResponder(http.StatusCreated, `{
								  "transaction_id": "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
								  "status": "failed",
								  "code": "100"
								}`))
	httpmock.RegisterMatcherResponder("POST", "http://bank_simulator:8081/api/v1/deposit",
		httpmock.BodyContainsString(`"number":"4000000000009995"`),
		httpmock.NewStringResponder(http.StatusCreated, `{
								  "transaction_id": "f17c2a1c-2459-40c9-9f2b-7587491d7df8",
								  "status": "pending",
								  "code": "200"
								}`))

	// mock UUID generator
	mockUUID := mock.NewMockUUIDGenerator(ctrl)
	mockUUID.EXPECT().New().Return("mockID").AnyTimes()

	return payment.NewPaymentService(
		mockRepository,
		cardProcessor,
		mockCardRepository,
		mockUUID,
		nil,
	)
}
