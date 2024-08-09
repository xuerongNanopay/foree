package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	sqlEmailPasswdRecoverInsert = `
	INSERT INTO email_passwd_recover
		(
			code, email_passwd_id, expiry_at
		) VALUES(?,?,?)
	`
	sQLEmailPasswdRecoverUpdateByCode = `
		UPDATE email_passwd_recover SET
			is_redeemed = ?, redeem_at = ?
		WHERE code = ?
	`
	sQLEmailPasswdRecoverGetUniqueByCode = `
		SELECT 
			p.code, p.email_passwd_id, p.is_redeemed, p.redeem_at
			p.expiry_at, p.create_at, p.update_at
        FROM email_passwd_recover p
        where p.code = ?
	`
)

type EmailPasswdRecover struct {
	Code          string    `json:"code"`
	EmailPasswdId int64     `json:"emailPasswdId"`
	IsRedeemed    bool      `json:"isRedeemed"`
	RedeemAt      time.Time `json:"redeemAt"`
	ExpiryAt      time.Time `json:"expiryAt"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewEmailPasswdRecoverRepo(db *sql.DB) *EmailPasswdRecoverRepo {
	return &EmailPasswdRecoverRepo{db: db}
}

type EmailPasswdRecoverRepo struct {
	db *sql.DB
}

func (repo *EmailPasswdRecoverRepo) InsertEmailPasswdRecover(p EmailPasswdRecover) (string, error) {
	p.Code = generateEmailPasswdRecoverCode(p.EmailPasswdId)
	result, err := repo.db.Exec(
		sqlEmailPasswdRecoverInsert,
		p.Code,
		p.EmailPasswdId,
		p.ExpiryAt,
	)
	if err != nil {
		return "", err
	}
	_, qerr := result.LastInsertId()
	if qerr != nil {
		return "", qerr
	}
	return p.Code, nil
}

func (repo *EmailPasswdRecoverRepo) UpdateEmailPasswdRecoverById(p EmailPasswdRecover) error {
	_, err := repo.db.Exec(sQLEmailPasswdRecoverUpdateByCode, p.IsRedeemed, p.RedeemAt, p.Code)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRecoverRepo) GetUniqueEmailPasswdRecoverByCode(code string) (*EmailPasswdRecover, error) {
	rows, err := repo.db.Query(sQLEmailPasswdRecoverGetUniqueByCode, code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *EmailPasswdRecover

	for rows.Next() {
		f, err = scanRowIntoEmailPasswdRecover(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.Code == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoEmailPasswdRecover(rows *sql.Rows) (*EmailPasswdRecover, error) {
	u := new(EmailPasswdRecover)
	err := rows.Scan(
		&u.Code,
		&u.EmailPasswdId,
		&u.IsRedeemed,
		&u.RedeemAt,
		&u.ExpiryAt,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func generateEmailPasswdRecoverCode(emailPasswdId int64) string {
	return fmt.Sprintf("%012d-%s", emailPasswdId, uuid.New().String())
}
