package payment

import "context"

type Payment struct {
	ID         string
	MerchantID string
	TrackingID string // id of the payment in the payment processor
	CardInfo   CardInfo
	Status     string
	Amount     int64
	Currency   string
}

type CardInfo struct {
	ID          string
	CardToken   string
	ExpiryMonth string
	ExpiryYear  string
	CVV         string
	CardNumber  string // todo put card to a separate table, store mask
}

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
	Deposit(ctx context.Context, cardToken string, amount int64, currency string) error
}

type Service interface {
	ProcessPayment(payment *Payment) (*Payment, error)
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

func (s *ServiceImpl) ProcessPayment(payment *Payment) (*Payment, error) {
	// Use a payment processor to tokenize the card data and retrieve a card token
	//card := payment.CardInfo

	// todo tokenize card
	//_, err := s.pi.TokenizeCard(card.CardNumber, card.ExpiryMonth, card.ExpiryYear, card.CVV)
	//if err != nil {
	//	return nil, err
	//}

	// Create a payment entity with the card token and other payment details
	payment.ID = s.uuid.New()
	payment.Status = "initiated"

	// Store the payment in the repository
	err := s.repo.Store(payment)
	if err != nil {
		return nil, err
	}

	err = s.cardProcessor.Deposit(context.Background(), payment.CardInfo.CardToken, payment.Amount, payment.Currency)
	if err != nil {
		return nil, err
	}

	s.repo.Store(payment)

	// Return the payment entity with the masked card number
	//payment.CardNumber = maskCardNumber(cardNumber)
	return payment, nil
}

func (s *ServiceImpl) GetPaymentDetails(id string) (*Payment, error) {
	// Retrieve the payment entity from the repository
	payment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Mask the card number before returning the payment entity
	// todo mask card number
	return payment, nil
}
