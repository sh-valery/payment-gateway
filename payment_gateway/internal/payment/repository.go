package payment

import "database/sql"

type MysqlRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *MysqlRepositoryImpl {
	return &MysqlRepositoryImpl{db: db}
}

func (r *MysqlRepositoryImpl) Store(p *Payment) error {
	query := `insert into payments (uuid, merchant_id, tracking_id, card_token, status, status_code, amount, currency) 
		    values (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		p.ID, p.MerchantID, p.TrackingID, p.CardInfo.CardToken, p.Status, p.StatusCode, p.Amount, p.Currency)
	if err != nil {
		return err
	}

	return nil
}

func (r *MysqlRepositoryImpl) GetByID(ID string) (*Payment, error) {
	p := &Payment{}

	err := r.db.QueryRow(`
		SELECT uuid, merchant_id, tracking_id, card_token, status, amount, currency 
		FROM payments 
		WHERE uuid = $1`, ID).Scan(&p.ID, &p.MerchantID, &p.TrackingID, &p.CardInfo.CardToken, &p.Status, &p.Amount, &p.Currency)
	if err != nil {
		return nil, err
	}

	return p, nil
}
