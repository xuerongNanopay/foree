package auth

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

const (
	sQLEmailPasswdInsert = `
		INSERT INTO email_passwd
		(	u.email, u.password, u.status, u.verify_code, u.user_id
		) VALUES (?,?,?,?)
	`
	sQLEmailPasswdUpdateByEmail = `
		UPDATE email_passwd SET 
			status = ?, password = ?, verify_code = ?, code_verified_at = ?
		WHERE email = ?
	`
	sQLEmailPasswdGetUniqueById = `
	SELECT 
		u.id, u.email, u.password, u.status,
		u.verify_code, u.code_verified_at,
		u.user_id, u.create_at, u.update_at
	FROM email_passwd as u 
	WHERE u.id = ?
`
	sQLEmailPasswdGetUniqueByEmail = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.user_id, u.create_at, u.update_at
		FROM email_passwd as u 
		WHERE u.email = ?
	`
	sQLEmailPasswdGetAll = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.user_id, u.create_at, u.update_at
		FROM email_passwd as u
	`
)

type EmailPasswdStatus string

const (
	EPStatusWaitingVerify EmailPasswdStatus = "WAITING_VERIFY"
	EPStatusPassExpire    EmailPasswdStatus = "PASSWORD_EXPIRE"
	EPStatusActive        EmailPasswdStatus = "ACTIVE"
	EPStatusSuspend       EmailPasswdStatus = "SUSPEND"
	EPStatusDisable       EmailPasswdStatus = "DISABLE"
)

type EmailPasswd struct {
	ID             int64             `json:"id"`
	Status         EmailPasswdStatus `json:"status"`
	Email          string            `json:"email"`
	Passowrd       string            `json:"-"`
	VerifyCode     string            `json:"-"`
	CodeVerifiedAt time.Time         `json:"codeVerifiedAt"`
	UserId         int64             `json:"userId"`
	CreateAt       time.Time         `json:"createAt"`
	UpdateAt       time.Time         `json:"updateAt"`
}

func NewEmailPasswdRepo(db *sql.DB) *EmailPasswdRepo {
	return &EmailPasswdRepo{db: db}
}

type EmailPasswdRepo struct {
	db *sql.DB
}

func (repo *EmailPasswdRepo) InsertEmailPasswd(ep EmailPasswd) (int64, error) {
	result, err := repo.db.Exec(
		sQLEmailPasswdInsert,
		ep.Email,
		ep.Passowrd,
		ep.Status,
		ep.VerifyCode,
		ep.UserId,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
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

func (repo *EmailPasswdRepo) UpdateEmailPasswdByEmail(ep EmailPasswd) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdateByEmail, ep.Status, ep.Passowrd, ep.VerifyCode, ep.CodeVerifiedAt, ep.Email)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoEmailPasswd(rows *sql.Rows) (*EmailPasswd, error) {
	p := new(EmailPasswd)
	err := rows.Scan(
		&p.ID,
		&p.Email,
		&p.Passowrd,
		&p.Status,
		&p.VerifyCode,
		&p.CodeVerifiedAt,
		&p.UserId,
		&p.CreateAt,
		&p.UpdateAt,
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
func HashPassword(password string) (string, error) {
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// if err != nil {
	// 	return "", err
	// }

	// return string(hash), nil
	return password, nil
}

// TODO
func ComparePasswords(hashed string, plain []byte) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	// return err == nil
	return hashed == string(plain)
}
