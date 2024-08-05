package auth

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	SQLEmailPasswdInsert = `
		INSERT INTO email_passwd
		(	u.email, u.password, u.status, u.verify_code
		) VALUES (?,?,?,?)
	`
	SQLEmailPasswdGetUniqueByEmail = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.create_at, u.update_at
		FROM email_passwd as u 
		WHERE u.email = ?
	`
	SQLEmailPasswdGetAll = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.avatar_url, u.create_at, u.update_at
		FROM email_passwd as u
	`
	SQLEmailPasswdUpdateStatusByEmail = `
		UPDATE email_passwd SET status = ? WHERE email = ?
	`
	SQLEmailPasswdUpdatePasswdByEmail = `
		UPDATE email_passwd SET password = ? WHERE email = ?
	`
	SQLEmailPasswdUpdateVerifyCodeByEmail = `
		UPDATE email_passwd SET verify_code = ? WHERE email = ?
	`
	SQLEmailPasswdUpdateCodeVerifiedAtByEmail = `
		UPDATE email_passwd SET code_verified_at = ? WHERE email = ?
	`
	SQLEmailPasswdUpdateUserIdByEmail = `
		UPDATE email_passwd SET user_id = ? WHERE email = ?
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
	CreateAt       time.Time         `json:"createAt"`
	UpdateAt       time.Time         `json:"updateAt"`
	UserId         int64             `json:"userId"`
}

func NewEmailPasswdRepo(db *sql.DB) *EmailPasswdRepo {
	return &EmailPasswdRepo{db: db}
}

type EmailPasswdRepo struct {
	db *sql.DB
}

func (repo *EmailPasswdRepo) Insert(ep EmailPasswd) (int64, error) {
	result, err := repo.db.Exec(
		SQLEmailPasswdInsert,
		ep.Email,
		ep.Passowrd,
		ep.Status,
		ep.VerifyCode,
	)
	if err != nil {
		return 0, fmt.Errorf("Insert: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Insert: %v", err)
	}
	return id, nil
}

func (repo *EmailPasswdRepo) GetUniqueByEmail(email string) (*EmailPasswd, error) {
	rows, err := repo.db.Query(SQLEmailPasswdGetUniqueByEmail, email)

	if err != nil {
		return nil, fmt.Errorf("GetUniqueByEmail: %v", err)
	}
	defer rows.Close()

	var ep *EmailPasswd

	for rows.Next() {
		ep, err = scanRowIntoEmailPasswd(rows)
		if err != nil {
			return nil, fmt.Errorf("GetUniqueByEmail: %v", err)
		}
	}

	if ep.ID == 0 {
		return nil, fmt.Errorf("GetUniqueByEmail: id `%v` not found", email)
	}

	return ep, nil
}

func (repo *EmailPasswdRepo) GetAllByEmail() ([]*EmailPasswd, error) {
	rows, err := repo.db.Query(SQLEmailPasswdGetAll)

	if err != nil {
		return nil, fmt.Errorf("GetAll: %v", err)
	}
	defer rows.Close()

	var eps []*EmailPasswd

	for rows.Next() {
		ep, err := scanRowIntoEmailPasswd(rows)
		if err != nil {
			return nil, fmt.Errorf("GetAll: %v", err)
		}
		eps = append(eps, ep)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAll: %v", err)
	}

	return eps, nil
}

func (repo *EmailPasswdRepo) UpdateStatusByEmail(email string, status EmailPasswdStatus) error {
	_, err := repo.db.Exec(SQLEmailPasswdUpdateStatusByEmail, status, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdatePasswdByEmail(email string, passwd string) error {
	_, err := repo.db.Exec(SQLEmailPasswdUpdatePasswdByEmail, passwd, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateVerifyCodeByEmail(email string, newCode string) error {
	_, err := repo.db.Exec(SQLEmailPasswdUpdateVerifyCodeByEmail, newCode, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateCodeVerifiedAtByEmail(email string, t time.Time) error {
	_, err := repo.db.Exec(SQLEmailPasswdUpdateCodeVerifiedAtByEmail, t, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateUserIdByEmail(email string, userId int64) error {
	_, err := repo.db.Exec(SQLEmailPasswdUpdateUserIdByEmail, userId, email)
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
		&p.CreateAt,
		&p.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
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
