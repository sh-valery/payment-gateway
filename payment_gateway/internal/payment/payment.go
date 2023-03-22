package payment

import (
	"context"
	"errors"
	"log"
)

//go:generate mockgen -source=payment.go -destination=./mock/payment.go -package=mock
type Repository interface {
	Create(payment *Payment) error
	Update(id, status, statusCode, trackingID string) error
	GetByID(id string) (*Payment, error)
}

// CardRepository should be PCI DSS repository
type CardRepository interface {
	SaveCardInfo(card *CardInfo) error
	GetCardByID(id string) (*CardInfo, error)
}

type UUIDGenerator interface {
	New() string
}

type CardProcessor interface {
	Deposit(ctx context.Context, payment *Payment) error
}

type Service interface {
	ProcessPayment(payment *Payment) error
	GetPaymentDetails(id string, merchantID string) (*Payment, error)
}

type ServiceImpl struct {
	repo          Repository
	cardProcessor CardProcessor
	pi            CardRepository
	uuid          UUIDGenerator
	// infra layer
	logger *log.Logger
}

func NewPaymentService(repo Repository, processor CardProcessor, tokenizer CardRepository, uuid UUIDGenerator, logger *log.Logger) Service {
	if logger == nil {
		logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	}

	return &ServiceImpl{
		repo:          repo,
		cardProcessor: processor,
		pi:            tokenizer,
		uuid:          uuid,
		logger:        logger,
	}
}

func (s *ServiceImpl) ProcessPayment(payment *Payment) error {
	// Save the card info in the PCI DSS repository
	card := payment.CardInfo
	card.ID = s.uuid.New()
	err := s.pi.SaveCardInfo(card)
	if err != nil {
		s.logger.Printf("error saving card info: %s", err)
		return err
	}
	s.logger.Printf("card %s info saved correctly: %s", card.ID, card.GetMaskedNumber())

	// Create a payment entity with the card token and other payment details
	payment.ID = s.uuid.New()
	payment.Status = "initiated"

	// Create the payment in the repository
	err = s.repo.Create(payment)
	if err != nil {
		s.logger.Printf("error creating payment: %s", err)
		return err
	}
	s.logger.Printf("payment %s init correctly", payment.ID)

	err = s.cardProcessor.Deposit(context.Background(), payment)
	if err != nil {
		s.logger.Printf("error depositing payment: %s", err)
		return err
	}
	s.logger.Printf("get resp from the card processor %+v", payment)

	err = s.repo.Update(payment.ID, payment.Status, payment.StatusCode, payment.TrackingID)
	if err != nil {
		s.logger.Printf("error updating payment: %s", err)
		return err
	}
	s.logger.Printf("payment: %s, trackingID: %s updated correctly to %s", payment.ID, payment.TrackingID, payment.Status)

	return nil
}

func (s *ServiceImpl) GetPaymentDetails(id string, merchantID string) (*Payment, error) {
	// Retrieve the payment entity from the repository
	payment, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Printf("error getting payment: %s", err)
		return nil, err
	}
	s.logger.Printf("payment %s retrieved correctly", payment.ID)

	if payment.MerchantID != merchantID {
		s.logger.Printf("merchant id does not match, expected: %s, got: %s", payment.MerchantID, merchantID)
		return nil, errors.New("merchant id does not match")
	}

	card, err := s.pi.GetCardByID(payment.CardInfo.ID)
	if err != nil {
		s.logger.Printf("error getting card info: %s", err)
		return nil, err
	}
	payment.CardInfo = card
	s.logger.Printf("card %s with number %s retrieved correctly", card.ID, card.GetMaskedNumber())

	return payment, nil
}
