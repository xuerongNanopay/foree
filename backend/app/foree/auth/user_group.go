package auth

import "time"

const (
	sQLUserGroupInsert = `
		INSERT INTO user_group
		(
			role_group, transaction_limit_group
		) VALUES(?,?,?)
	`
	sQLUserGoup
)

type UserGroup struct {
	ID                    int64     `json:"id"`
	RoleGroup             string    `json:"roleGroup"`
	TransactionLimitGroup string    `json:"transactionLimitGroup"`
	OwnerId               int64     `json:"ownerId"`
	CreateAt              time.Time `json:"createAt"`
	UpdateAt              time.Time `json:"updateAt"`
}
