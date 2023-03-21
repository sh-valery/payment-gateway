package payment

import "context"

type Repository interface {
	Store(payment *Payment) error
	GetByID(ID string) (*Payment, error)
}

type CardTokenizer interface {
	TokenizeCard(cardNumber, expiryMonth, expiryYear, cvv string) (string, error)
}

type UUIDGenerator interface {
	New() string
}

type CardProcessor interface {
	Deposit(ctx context.Context, payment *Payment) error
}

type Service interface {
	ProcessPayment(payment *Payment) error
	GetPaymentDetails(id string) (*Payment, error)
}

type ServiceImpl struct {
	repo          Repository
	cardProcessor CardProcessor
	pi            CardTokenizer
	uuid          UUIDGenerator
}

func NewPaymentService(repo Repository, processor CardProcessor, tokenizer CardTokenizer, uuid UUIDGenerator) Service {
	return &ServiceImpl{
		repo:          repo,
		cardProcessor: processor,
		pi:            tokenizer,
		uuid:          uuid,
	}
}

func (s *ServiceImpl) ProcessPayment(payment *Payment) error {
	// Use a payment processor to tokenize the card data and retrieve a card token
	//card := payment.CardInfo

	// todo tokenize card
	//_, err := s.pi.TokenizeCard(card.cardNumber, card.ExpiryMonth, card.ExpiryYear, card.cvv)
	//if err != nil {
	//	return nil, err
	//}

	// Create a payment entity with the card token and other payment details
	payment.ID = s.uuid.New()
	payment.Status = "initiated"

	// Store the payment in the repository
	err := s.repo.Store(payment)
	if err != nil {
		return err
	}

	err = s.cardProcessor.Deposit(context.Background(), payment)
	if err != nil {
		return err
	}

	err = s.repo.Store(payment)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetPaymentDetails(id string) (*Payment, error) {
	// Retrieve the payment entity from the repository
	payment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
