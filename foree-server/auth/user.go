package auth

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLUserInsert = `
		INSERT INTO users(	
			status, first_name, middle_name, 
			last_name, age, dob, 
			address1, address2, city, province, country, postal_code, 
			phone_number, email
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLUserUpdateById = `
		UPDATE users SET 
			status = ?, first_name = ?, middle_name = ?, 
			last_name = ?, age = ?, dob = ?, address1 = ?, 
			address2 = ?, city = ?, province = ?, country = ?, postal_code = ?, phone_number = ?,
			email = ?
		WHERE id = ?
	`
	sQLUserGetAll = `
		SELECT 
			u.id, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.address1, 
			u.address2, u.city, u.province, u.country, u.postal_code, u.phone_number,
			u.email, u.avatar_url, u.created_at, u.updated_at
		FROM users as u 
	`
	sQLUserGetUniqueById = `
		SELECT 
			u.id, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.address1, 
			u.address2, u.city, u.province, u.country, u.postal_code, u.phone_number,
			u.email, u.avatar_url, u.created_at, u.updated_at
		FROM users as u 
		WHERE u.id = ?
	`
)

type UserStatus string

const (
	UserStatusInitial UserStatus = "INITIAL"
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusSuspend UserStatus = "SUSPEND"
	UserStatusDelete  UserStatus = "DELETE"
)

type User struct {
	ID          int64      `json:"id"`
	Status      UserStatus `json:"status"`
	FirstName   string     `json:"firstName"`
	MiddleName  string     `json:"middleName"`
	LastName    string     `json:"lastName"`
	Age         int        `json:"age"`
	Dob         time.Time  `json:"dob"`
	Address1    string     `json:"address1"`
	Address2    string     `json:"address2"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Country     string     `json:"country"`
	PostalCode  string     `json:"postalCode"`
	PhoneNumber string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	AvatarUrl   string     `json:"avatarUrl"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

type UserRepo struct {
	db *sql.DB
}

func (repo *UserRepo) InsertUser(ctx context.Context, user User) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLUserInsert,
			user.Status,
			user.FirstName,
			user.MiddleName,
			user.LastName,
			user.Age,
			user.Dob,
			user.Address1,
			user.Address2,
			user.City,
			user.Province,
			user.Country,
			user.PostalCode,
			user.PhoneNumber,
			user.Email,
		)
	} else {
		result, err = repo.db.Exec(
			sQLUserInsert,
			user.Status,
			user.FirstName,
			user.MiddleName,
			user.LastName,
			user.Age,
			user.Dob,
			user.Address1,
			user.Address2,
			user.City,
			user.Province,
			user.Country,
			user.PostalCode,
			user.PhoneNumber,
			user.Email,
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

func (repo *UserRepo) UpdateUserById(ctx context.Context, u User) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLUserUpdateById,
			u.Status,
			u.FirstName,
			u.MiddleName,
			u.LastName,
			u.Age,
			u.Dob,
			u.Address1,
			u.Address2,
			u.City,
			u.Province,
			u.Country,
			u.PostalCode,
			u.PhoneNumber,
			u.Email,
			u.ID,
		)
	} else {
		_, err = repo.db.Exec(
			sQLUserUpdateById,
			u.Status,
			u.FirstName,
			u.MiddleName,
			u.LastName,
			u.Age,
			u.Dob,
			u.Address1,
			u.Address2,
			u.City,
			u.Province,
			u.Country,
			u.PostalCode,
			u.PhoneNumber,
			u.Email,
			u.ID,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) GetUniqueUserById(id int64) (*User, error) {
	rows, err := repo.db.Query(sQLUserGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *User

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u == nil || u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func (repo *UserRepo) GetAllUser() ([]*User, error) {
	rows, err := repo.db.Query(sQLUserGetAll)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		u, err := scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func scanRowIntoUser(rows *sql.Rows) (*User, error) {
	u := new(User)
	err := rows.Scan(
		&u.ID,
		&u.Status,
		&u.FirstName,
		&u.MiddleName,
		&u.LastName,
		&u.Age,
		&u.Dob,
		&u.Address1,
		&u.Address2,
		&u.City,
		&u.Province,
		&u.Country,
		&u.PostalCode,
		&u.PhoneNumber,
		&u.Email,
		&u.AvatarUrl,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
