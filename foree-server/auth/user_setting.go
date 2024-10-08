package auth

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLUserSettingInsert = `
		INSERT INTO user_setting(
			is_in_app_notification_enable, is_push_notification_enable, is_email_notifications_enable, owner_id
		) VALUES(?,?,?,?)
	`
	sQLUserSettingUpdate = `
		UPDATE user_setting SET
			is_in_app_notification_enable = ?, is_push_notification_enable = ?, is_email_notifications_enable = ?
		WHERE owner_id = ?
	`
	sQLUserSettingUniqueByOwnerId = `
		SELECT
			u.id, u.is_in_app_notification_enable, u.is_push_notification_enable, u.is_email_notifications_enable,
			u.owner_id, u.created_at, u.updated_at
		FROM user_setting as u
		WHERE u.owner_id = ?
	`
)

type UserSetting struct {
	ID                         int64      `json:"id"`
	IsInAppNotificationEnable  bool       `json:"isInAppNotificationEnable"`
	IsPushNotificationEnable   bool       `json:"isPushNotificationEnable"`
	IsEmailNotificationsEnable bool       `json:"isEmailNotificationsEnable"`
	OwnerId                    int64      `json:"ownerId"`
	CreatedAt                  *time.Time `json:"createdAt"`
	UpdatedAt                  *time.Time `json:"updatedAt"`
}

type UserSettingRepo struct {
	db *sql.DB
}

func NewUserSettingRepo(db *sql.DB) *UserSettingRepo {
	return &UserSettingRepo{db: db}
}

func (repo *UserSettingRepo) InsertUserSetting(ctx context.Context, us UserSetting) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLUserSettingInsert,
			us.IsInAppNotificationEnable,
			us.IsPushNotificationEnable,
			us.IsEmailNotificationsEnable,
			us.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLUserSettingInsert,
			us.IsInAppNotificationEnable,
			us.IsPushNotificationEnable,
			us.IsEmailNotificationsEnable,
			us.OwnerId,
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

func (repo *UserSettingRepo) UpdateUserSettingByOwnerId(ctx context.Context, us UserSetting) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLUserSettingUpdate,
			us.IsInAppNotificationEnable,
			us.IsPushNotificationEnable,
			us.IsEmailNotificationsEnable,
			us.OwnerId,
		)
	} else {
		_, err = repo.db.Exec(
			sQLUserSettingUpdate,
			us.IsInAppNotificationEnable,
			us.IsPushNotificationEnable,
			us.IsEmailNotificationsEnable,
			us.OwnerId,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *UserSettingRepo) GetUniqueUserSettingByOwnerId(ownerId int64) (*UserSetting, error) {
	rows, err := repo.db.Query(sQLUserSettingUniqueByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *UserSetting

	for rows.Next() {
		u, err = scanRowUserSetting(rows)
		if err != nil {
			return nil, err
		}
	}

	if u == nil || u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func scanRowUserSetting(rows *sql.Rows) (*UserSetting, error) {
	u := new(UserSetting)
	err := rows.Scan(
		&u.ID,
		&u.IsInAppNotificationEnable,
		&u.IsPushNotificationEnable,
		&u.IsEmailNotificationsEnable,
		&u.OwnerId,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
