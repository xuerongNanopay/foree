package auth

import (
	"database/sql"
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
	ID             uint64            `json:"id"`
	Status         EmailPasswdStatus `json:"status"`
	Email          string            `json:"email"`
	Passowrd       string            `json:"-"`
	VerifyCode     string            `json:"-"`
	CodeVerifiedAt time.Time         `json:"codeVerifiedAt"`
	CreateAt       time.Time         `json:"createAt"`
	UpdateAt       time.Time         `json:"updateAt"`
	UserId         uint64            `json:"userId"`
}

func NewEmailPasswdRepo(db *sql.DB) *EmailPasswdRepo {
	return &EmailPasswdRepo{db: db}
}

type EmailPasswdRepo struct {
	db *sql.DB
}

func (repo *EmailPasswdRepo) Insert(ep EmailPasswd) (*EmailPasswd, error) {
	return nil, nil
}

func (repo *EmailPasswdRepo) GetUniqueByEmail(email string) (*EmailPasswd, error) {
	return nil, nil
}

func (repo *EmailPasswdRepo) UpdateStatusByEmail(email string, status EmailPasswdStatus) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdatePasswdByEmail(email string, passwd string) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdateVerifyCodeByEmail(email string, newCode string) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdateCodeVerifiedAtByEmail(email string, t time.Time) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdateUserIdByEmail(email string, userId int64) error {
	return nil
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
