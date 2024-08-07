package auth

import (
	"database/sql"
	"time"
)

const (
	sQLEmailPasswdInsert = `
		INSERT INTO email_passwd
		(	u.email, u.password, u.status, u.verify_code
		) VALUES (?,?,?,?)
	`
	sQLEmailPasswdGetUniqueByEmail = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.create_at, u.update_at
		FROM email_passwd as u 
		WHERE u.email = ?
	`
	sQLEmailPasswdGetAll = `
		SELECT 
			u.id, u.email, u.password, u.status,
			u.verify_code, u.code_verified_at,
			u.avatar_url, u.create_at, u.update_at
		FROM email_passwd as u
	`
	sQLEmailPasswdUpdateStatusByEmail = `
		UPDATE email_passwd SET status = ? WHERE email = ?
	`
	sQLEmailPasswdUpdatePasswdByEmail = `
		UPDATE email_passwd SET password = ? WHERE email = ?
	`
	sQLEmailPasswdUpdateVerifyCodeByEmail = `
		UPDATE email_passwd SET verify_code = ? WHERE email = ?
	`
	sQLEmailPasswdUpdateCodeVerifiedAtByEmail = `
		UPDATE email_passwd SET code_verified_at = ? WHERE email = ?
	`
	sQLEmailPasswdUpdateUserIdByEmail = `
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
		sQLEmailPasswdInsert,
		ep.Email,
		ep.Passowrd,
		ep.Status,
		ep.VerifyCode,
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

func (repo *EmailPasswdRepo) GetUniqueByEmail(email string) (*EmailPasswd, error) {
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

func (repo *EmailPasswdRepo) GetAllByEmail() ([]*EmailPasswd, error) {
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

func (repo *EmailPasswdRepo) UpdateStatusByEmail(email string, status EmailPasswdStatus) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdateStatusByEmail, status, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdatePasswdByEmail(email string, passwd string) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdatePasswdByEmail, passwd, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateVerifyCodeByEmail(email string, newCode string) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdateVerifyCodeByEmail, newCode, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateCodeVerifiedAtByEmail(email string, t time.Time) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdateCodeVerifiedAtByEmail, t, email)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EmailPasswdRepo) UpdateUserIdByEmail(email string, userId int64) error {
	_, err := repo.db.Exec(sQLEmailPasswdUpdateUserIdByEmail, userId, email)
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
