package payment

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
)

type MysqlRepositoryImpl struct {
	db *sql.DB

	cryptoKey []byte
}

func NewPaymentRepository(db *sql.DB) *MysqlRepositoryImpl {
	return &MysqlRepositoryImpl{
		db:        db,
		cryptoKey: []byte("mock_key_16b_len"), // todo put to config 16b
	}
}

func (r *MysqlRepositoryImpl) Create(p *Payment) error {
	// todo card_token -> card_id
	query := `insert into payments (uuid, merchant_id, tracking_id, card_token, status, status_code, amount, currency) 
		    values (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		p.ID, p.MerchantID, p.TrackingID, p.CardInfo.ID, p.Status, p.StatusCode, p.Amount, p.Currency)
	if err != nil {
		return err
	}

	return nil
}

func (r *MysqlRepositoryImpl) GetByID(ID string) (*Payment, error) {
	p := &Payment{
		CardInfo: &CardInfo{},
	}

	query := `SELECT uuid, merchant_id, tracking_id, card_token, status, status_code, amount, currency 
		FROM payments 
		WHERE uuid = ?`
	err := r.db.QueryRow(query, ID).
		Scan(&p.ID, &p.MerchantID, &p.TrackingID, &p.CardInfo.ID, &p.Status, &p.StatusCode, &p.Amount, &p.Currency)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *MysqlRepositoryImpl) UpdateStatusTracking(payment *Payment) error {
	query := `update payments set status = ?, status_code = ?, tracking_id = ? where uuid = ?`

	res, err := r.db.Exec(query, payment.Status, payment.StatusCode, payment.TrackingID, payment.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *MysqlRepositoryImpl) SaveCardInfo(card *CardInfo) error {
	encryptedCard, err := r.encrypt(card.cardNumber)
	encryptedCVV, err := r.encrypt(card.cvv)
	encryptedHolderName, err := r.encrypt(card.holderName)
	if err != nil {
		return errors.New("error encrypting card info")
	}

	query := `insert into cards (uuid, card_number, card_holder, expiry_month, expiry_year, cvv)
    		values (?, ?, ?, ?, ?, ?)`
	_, err = r.db.Exec(query, card.ID, encryptedCard, encryptedHolderName, card.ExpiryMonth, card.ExpiryYear, encryptedCVV)
	if err != nil {
		return err
	}
	return nil
}

func (r *MysqlRepositoryImpl) GetCardByID(ID string) (*CardInfo, error) {
	var encryptedCard, encryptedCVV, encryptedHolderName string
	query := `SELECT uuid, card_number, card_holder, expiry_month, expiry_year, cvv
		FROM cards 
		where uuid = ?`
	card := &CardInfo{}
	err := r.db.QueryRow(query, ID).
		Scan(&card.ID, &encryptedCard, &encryptedHolderName, &card.ExpiryMonth, &card.ExpiryYear, &encryptedCVV)
	if err != nil {
		return nil, err
	}

	card.cardNumber, err = r.decrypt(encryptedCard)
	card.holderName, err = r.decrypt(encryptedHolderName)
	card.cvv, err = r.decrypt(encryptedCVV)
	if err != nil {
		return nil, errors.New("error decrypting card info")
	}

	return card, nil
}

func (r *MysqlRepositoryImpl) encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(r.cryptoKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	encryptedText := base64.URLEncoding.EncodeToString(ciphertext)
	return encryptedText, nil
}

func (r *MysqlRepositoryImpl) decrypt(text string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(r.cryptoKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
