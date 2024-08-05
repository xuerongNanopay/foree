package auth

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	SQLGetUserByEmail = `
		SELECT 
			u.id, u.group, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.nationality, u.Address1, 
			u.Address2, u.city, u.province, u.country, u.phone_number,
			u.email, u.avatar_url, u.create_at, u.update_at
		FROM users as u WHERE u.email = ?
	`
	SQLGetUserById = `
		SELECT 
			u.id, u.group, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.nationality, u.Address1, 
			u.Address2, u.city, u.province, u.country, u.phone_number,
			u.email, u.avatar_url, u.create_at, u.update_at
		FROM users as u WHERE u.id = ?
	`
	SQLInsertUser = `
		INSERT INTO users
		(	id, group, status, first_name, middle_name, 
			last_name, age, dob, nationality, Address1, 
			Address2, city, province, country, phone_number,
			email
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
	`
	SQLUpdateUserStatus = `
		UPDATE user SET status = ? WHERE id = ?
	`
)

type UserStatus string

const (
	UserStatusInitial UserStatus = "INITIAL"
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusSuspend UserStatus = "SUSPEND"
	UserStatusDisable UserStatus = "DISABLE"
)

type User struct {
	ID          uint64     `json:"id"`
	Group       string     `json:"group"`
	Status      UserStatus `json:"status"`
	FirstName   string     `json:"firstName"`
	MiddleName  string     `json:"middleName"`
	LastName    string     `json:"lastName"`
	Age         int        `json:"age"`
	Dob         time.Time  `json:"dob"`
	Nationality string     `json:"nationality"`
	Address1    string     `json:"address1"`
	Address2    string     `json:"address2"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Country     string     `json:"country"`
	PhoneNumber string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	AvatarUrl   string     `json:"avatarUrl"`
	CreateAt    time.Time  `json:"createAt"`
	UpdateAt    time.Time  `json:"updateAt"`
	// OccupationId int64      `json:"-"`
	// Occupation   string     `json:"occupation"`
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

type UserRepo struct {
	db *sql.DB
}

func (repo *UserRepo) UpdateUserStatus(userId int64, status UserStatus) error {
	_, err := repo.db.Exec(SQLUpdateUserStatus, status, userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) InsertUser(user User) (int64, error) {
	result, err := repo.db.Exec(SQLInsertUser)
	if err != nil {
		return 0, fmt.Errorf("InsertUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertUser: %v", err)
	}
	return id, nil
}

func (repo *UserRepo) GetUserById(id int64) (*User, error) {
	rows, err := repo.db.Query(SQLGetUserById, id)

	if err != nil {
		return nil, fmt.Errorf("GetUserById: %v", err)
	}

	var u *User

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("GetUserById: %v", err)
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("GetUserByEmail: id `%v` not found", id)
	}

	return u, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*User, error) {
	rows, err := repo.db.Query(SQLGetUserByEmail, email)

	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %v", err)
	}

	var u *User

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, fmt.Errorf("GetUserByEmail: %v", err)
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("GetUserByEmail: email `%v` not found", email)
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*User, error) {
	u := new(User)
	err := rows.Scan(
		&u.ID,
		&u.Group,
		&u.Status,
		&u.FirstName,
		&u.MiddleName,
		&u.LastName,
		&u.Age,
		&u.Dob,
		&u.Nationality,
		&u.Address1,
		&u.Address2,
		&u.City,
		&u.Province,
		&u.Country,
		&u.PhoneNumber,
		&u.Email,
		&u.AvatarUrl,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
