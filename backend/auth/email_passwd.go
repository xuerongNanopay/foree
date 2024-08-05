package auth

import (
	"database/sql"
	"time"
)

const (
	SQLEmailPasswdGetUniqueByEmail = `
		SELECT 
			u.id, u.status, u.email, u.password, 
			u.verify_code, u.code_verified_at,
			u.avatar_url, u.create_at, u.update_at
		FROM email_passwd as u 
		WHERE u.email = ?
	`
	SQLEmailPasswdGetAll = `
		SELECT 
			u.id, u.status, u.email, u.password, 
			u.verify_code, u.code_verified_at,
			u.avatar_url, u.create_at, u.update_at
		FROM email_passwd as u
	`
	SQLEmailPasswdInsert = `
	INSERT INTO users
	(	id, group, status, first_name, middle_name, 
		last_name, age, dob, nationality, Address1, 
		Address2, city, province, country, phone_number,
		email
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
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

func (repo *EmailPasswdRepo) UpdateStatus(id int64, status EmailPasswdStatus) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdatePassword(id int64, passwd string) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdateVerifyCode(id int64, newCode string) error {
	return nil
}

func (repo *EmailPasswdRepo) UpdateUserId(id int64, userId int64) error {
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
