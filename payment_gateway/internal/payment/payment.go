package payment

import (
	"context"
	"errors"
)

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
}

func NewPaymentService(repo Repository, processor CardProcessor, tokenizer CardRepository, uuid UUIDGenerator) Service {
	return &ServiceImpl{
		repo:          repo,
		cardProcessor: processor,
		pi:            tokenizer,
		uuid:          uuid,
	}
}

func (s *ServiceImpl) ProcessPayment(payment *Payment) error {
	// Save the card info in the PCI DSS repository
	card := payment.CardInfo
	card.ID = s.uuid.New()
	err := s.pi.SaveCardInfo(card)
	if err != nil {
		return err
	}
	// Create a payment entity with the card token and other payment details
	payment.ID = s.uuid.New()
	payment.Status = "initiated"

	// Create the payment in the repository
	err = s.repo.Create(payment)
	if err != nil {
		return err
	}

	err = s.cardProcessor.Deposit(context.Background(), payment)
	if err != nil {
		return err
	}

	err = s.repo.Update(payment.ID, payment.Status, payment.StatusCode, payment.TrackingID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetPaymentDetails(id string, merchantID string) (*Payment, error) {
	// Retrieve the payment entity from the repository
	payment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if payment.MerchantID != merchantID {
		return nil, errors.New("merchant id does not match")
	}

	card, err := s.pi.GetCardByID(payment.CardInfo.ID)
	if err != nil {
		return nil, err
	}
	payment.CardInfo = card

	return payment, nil
}
