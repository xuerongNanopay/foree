package auth

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLEmailPasswdInsert = `
		INSERT INTO email_passwd(	 
			u.status, u.email, u.username, u.passwd, u.verify_code, u.owner_id
		) VALUES (?,?,?,?,?)
	`
	sQLEmailPasswdUpdateByEmail = `
		UPDATE email_passwd SET 
			status = ?, passwd = ?, verify_code = ?, verify_code_expired_at = ?, login_attempts = ?, retrieve_token = ?, retrieve_token_expired_at = ?
		WHERE email = ?
	`
	sQLEmailPasswdGetUniqueById = `
		SELECT 
			u.id, u.email, u.username, u.passwd, u.status,
			u.verify_code, u.verify_code_expired_at, u.login_attempts,
			u.retrieve_token, u.retrieve_token_expired_at,
			u.owner_id, u.created_at, u.updated_at
		FROM email_passwd as u 
		WHERE u.id = ?
`
	sQLEmailPasswdGetUniqueByEmail = `
		SELECT 
			u.id, u.email, u.username, u.passwd, u.status,
			u.verify_code, u.verify_code_expired_at, u.login_attempts,
			u.retrieve_token, u.retrieve_token_expired_at,
			u.owner_id, u.created_at, u.updated_at
		FROM email_passwd as u 
		WHERE u.email = ?
	`
	sQLEmailPasswdGetAll = `
		SELECT 
			u.id, u.email, u.username, u.passwd, u.status,
			u.verify_code, u.verify_code_expired_at, u.login_attempts,
			u.retrieve_token, u.retrieve_token_expired_at,
			u.owner_id, u.created_at, u.updated_at
		FROM email_passwd as u
	`
)

type EmailPasswdStatus string

const (
	EPStatusWaitingVerify EmailPasswdStatus = "WAITING_VERIFY"
	EPStatusPassExpire    EmailPasswdStatus = "PASSWORD_EXPIRE"
	EPStatusActive        EmailPasswdStatus = "ACTIVE"
	EPStatusSuspend       EmailPasswdStatus = "SUSPEND"
	EPStatusDelete        EmailPasswdStatus = "DELETE"
)

type EmailPasswd struct {
	ID                     int64             `json:"id"`
	Status                 EmailPasswdStatus `json:"status"`
	Email                  string            `json:"email"`
	Username               string            `json:"username"`
	Passwd                 string            `json:"-"`
	VerifyCode             string            `json:"-"`
	VerifyCodeExpiredAt    time.Time         `json:"verifyCodeExpiredAt"`
	RetrieveToken          string            `json:"-"`
	RetrieveTokenExpiredAt time.Time         `json:"retrieveTokenExpiredAt"`
	LoginAttempts          int32             `json:"loginAttempts"`
	OwnerId                int64             `json:"ownerId"`
	CreatedAt              time.Time         `json:"createdAt"`
	UpdatedAt              time.Time         `json:"updatedAt"`
}

func NewEmailPasswdRepo(db *sql.DB) *EmailPasswdRepo {
	return &EmailPasswdRepo{db: db}
}

type EmailPasswdRepo struct {
	db *sql.DB
}

func (repo *EmailPasswdRepo) InsertEmailPasswd(ctx context.Context, ep EmailPasswd) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLEmailPasswdInsert,
			ep.Status,
			ep.Email,
			ep.Username,
			ep.Passwd,
			ep.VerifyCode,
			ep.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLEmailPasswdInsert,
			ep.Status,
			ep.Email,
			ep.Username,
			ep.Passwd,
			ep.VerifyCode,
			ep.OwnerId,
		)
	}

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *EmailPasswdRepo) UpdateEmailPasswdByEmail(ctx context.Context, ep EmailPasswd) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLEmailPasswdUpdateByEmail,
			ep.Status,
			ep.Passwd,
			ep.VerifyCode,
			ep.VerifyCodeExpiredAt,
			ep.LoginAttempts,
			ep.RetrieveToken,
			ep.RetrieveTokenExpiredAt,
			ep.Email,
		)
	} else {
		_, err = repo.db.Exec(
			sQLEmailPasswdUpdateByEmail,
			ep.Status,
			ep.Passwd,
			ep.VerifyCode,
			ep.VerifyCodeExpiredAt,
			ep.LoginAttempts,
			ep.RetrieveToken,
			ep.RetrieveTokenExpiredAt,
			ep.Email,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) GetUniqueEmailPasswdByEmail(email string) (*EmailPasswd, error) {
	rows, err := repo.db.Query(sQLEmailPasswdGetUniqueByEmail, email)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ep *EmailPasswd

	for rows.Next() {
		ep, err = scanRowIntoEmailPasswd(rows)
		if err != nil {
			return nil, err
		}
	}

	if ep.ID == 0 {
		return nil, nil
	}

	return ep, nil
}

func (repo *EmailPasswdRepo) GetUniqueEmailPasswdById(id int64) (*EmailPasswd, error) {
	rows, err := repo.db.Query(sQLEmailPasswdGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ep *EmailPasswd

	for rows.Next() {
		ep, err = scanRowIntoEmailPasswd(rows)
		if err != nil {
			return nil, err
		}
	}

	if ep.ID == 0 {
		return nil, nil
	}

	return ep, nil
}

func (repo *EmailPasswdRepo) GetAllEmailPasswdByEmail() ([]*EmailPasswd, error) {
	rows, err := repo.db.Query(sQLEmailPasswdGetAll)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eps []*EmailPasswd

	for rows.Next() {
		ep, err := scanRowIntoEmailPasswd(rows)
		if err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return eps, nil
}

func scanRowIntoEmailPasswd(rows *sql.Rows) (*EmailPasswd, error) {
	p := new(EmailPasswd)
	err := rows.Scan(
		&p.ID,
		&p.Status,
		&p.Email,
		&p.Username,
		&p.Passwd,
		&p.VerifyCode,
		&p.VerifyCodeExpiredAt,
		&p.RetrieveToken,
		&p.RetrieveTokenExpiredAt,
		&p.OwnerId,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GenerateVerifyCode() string {
	r := rand.Intn(1000000)
	return fmt.Sprintf("%06d", r)
}

// TODO
func HashPassword(passwd string) (string, error) {
	// hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)

	// if err != nil {
	// 	return "", err
	// }

	// return string(hash), nil
	return passwd, nil
}

// TODO
func ComparePasswords(hashed string, plain []byte) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	// return err == nil
	return hashed == string(plain)
}
