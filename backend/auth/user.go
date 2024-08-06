package auth

import (
	"database/sql"
	"time"
)

const (
	SQLUserGetAll = `
		SELECT 
			u.id, u.group, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.nationality, u.Address1, 
			u.Address2, u.city, u.province, u.country, u.phone_number,
			u.email, u.avatar_url, u.create_at, u.update_at
		FROM users as u 
	`
	SQLUserGetUniqueById = `
		SELECT 
			u.id, u.group, u.status, u.first_name, u.middle_name, 
			u.last_name, u.age, u.dob, u.nationality, u.Address1, 
			u.Address2, u.city, u.province, u.country, u.phone_number,
			u.email, u.avatar_url, u.create_at, u.update_at
		FROM users as u 
		WHERE u.id = ?
	`
	SQLUserInsert = `
		INSERT INTO users
		(	group, status, first_name, middle_name, 
			last_name, age, dob, nationality, Address1, 
			Address2, city, province, country, phone_number,
			email
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	SQLUserUpdateStatus = `
		UPDATE users SET status = ? WHERE id = ?
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

func (repo *UserRepo) UpdateStatus(id int64, status UserStatus) error {
	_, err := repo.db.Exec(SQLUserUpdateStatus, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) Insert(user User) (int64, error) {
	result, err := repo.db.Exec(
		SQLUserInsert,
		user.Group,
		user.Status,
		user.FirstName,
		user.MiddleName,
		user.LastName,
		user.Age,
		user.Dob,
		user.Nationality,
		user.Address1,
		user.Address2,
		user.City,
		user.Province,
		user.Country,
		user.PhoneNumber,
		user.Email,
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

func (repo *UserRepo) GetUniqueById(id int64) (*User, error) {
	rows, err := repo.db.Query(SQLUserGetUniqueById, id)

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

	if u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func (repo *UserRepo) GetAll() ([]*User, error) {
	rows, err := repo.db.Query(SQLUserGetAll)

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
