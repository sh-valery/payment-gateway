package payment

import "database/sql"

type MysqlRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *MysqlRepositoryImpl {
	return &MysqlRepositoryImpl{db: db}
}

func (r *MysqlRepositoryImpl) Store(p *Payment) error {
	err := r.db.QueryRow(`
		INSERT INTO payments (id, merchant_id, tracking_id, card_token, status, amount, currency) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`,
		p.ID, p.MerchantID, p.TrackingID, p.CardInfo.CardToken, p.Status, p.Amount, p.Currency).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MysqlRepositoryImpl) GetByID(ID string) (*Payment, error) {
	p := &Payment{}

	err := r.db.QueryRow(`
		SELECT id, merchant_id, tracking_id, card_token, status, amount, currency 
		FROM payments 
		WHERE id = $1`, ID).Scan(&p.ID, &p.MerchantID, &p.TrackingID, &p.CardInfo.CardToken, &p.Status, &p.Amount, &p.Currency)
	if err != nil {
		return nil, err
	}

	return p, nil
}
