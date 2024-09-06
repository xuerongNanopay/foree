package promotion

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
)

const (
	sQLPromoCodeGetUniqueByCode = `
		SELECT
			P.code, p.description, p.min_amount, p.min_currency, p.limit_per_acc,
			p.is_enable, p.start_time, p.end_time, p.created_at, p.updated_at
		FROM promo_codes as p
		Where p.code = ?
	`
)

type PromoCode struct {
	Code        string           `json:"code"`
	Description string           `json:"description"`
	MinAmt      types.AmountData `json:"minAmt"`
	LimitPerAcc int              `json:"limit_per_acc"`
	IsEnable    bool             `json:"isEnable"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

func (p *PromoCode) IsValid() bool {
	if !p.IsEnable {
		return false
	}

	if now := time.Now().Unix(); now > p.StartTime.Unix() || (now > p.EndTime.Unix() && !p.EndTime.IsZero()) {
		return false
	}
	return true
}

func NewPromoCodeRepo(db *sql.DB) *PromoCodeRepo {
	return &PromoCodeRepo{db: db}
}

type PromoCodeRepo struct {
	db *sql.DB
}

func (repo *PromoCodeRepo) GetUniquePromoCodeByCode(ctx context.Context, code string) (*PromoCode, error) {
	rows, err := repo.db.Query(sQLPromoCodeGetUniqueByCode, code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *PromoCode

	for rows.Next() {
		f, err = scanRowIntoPromoCode(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.Code == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoPromoCode(rows *sql.Rows) (*PromoCode, error) {
	p := new(PromoCode)
	err := rows.Scan(
		&p.Code,
		&p.Description,
		&p.MinAmt.Amount,
		&p.MinAmt.Currency,
		&p.LimitPerAcc,
		&p.IsEnable,
		&p.StartTime,
		&p.EndTime,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

const (
	sQLPromoCodeJointInsert = `
		INSERT INTO promo_code_joint
		(
			status, promo_code, owner_id, transaction_id
		) VALUES(?,?,?,?)
	`
	sQLPromoCodeJointUpdateByTransactionId = `
		UPDATE promo_code_joint SET
			status = ?
		WHERE transaction_id = ?
	`
	sQLPromoCodeJointGetAllActiveByOwnerAndPromoCode = `
		SELECT
			j.id, j.status, j.promo_code, j.owner_id,
			j.transaction_id, f.created_at, f.updated_at
		FROM promo_code_joint as j
		WHERE j.owner_id = ? AND j.promo_code = ? AND j.status = INITIAL AND j.status = REDEEMED
	`
)

type PromoCodeJointStatus string

const (
	PromoCodeJointStatusInitial  PromoCodeJointStatus = "INITIAL"
	PromoCodeJointStatusRedeemed PromoCodeJointStatus = "REDEEMED"
	PromoCodeJointStatusDELETE   PromoCodeJointStatus = "DELETE"
)

type PromoCodeJoint struct {
	ID            int64                `json:"id"`
	Status        PromoCodeJointStatus `json:"status"`
	PromoCode     string               `json:"promoCode"`
	OwnerId       int64                `json:"ownerId"`
	TransactionId int64                `json:"transactionId"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

func NewPromoCodeJointRepo(db *sql.DB) *PromoCodeJointRepo {
	return &PromoCodeJointRepo{db: db}
}

type PromoCodeJointRepo struct {
	db *sql.DB
}

func (repo *PromoCodeJointRepo) InsertPromoCodeJoin(ctx context.Context, p PromoCodeJoint) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLPromoCodeJointInsert,
			p.Status,
			p.PromoCode,
			p.OwnerId,
			p.TransactionId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLPromoCodeJointInsert,
			p.Status,
			p.PromoCode,
			p.OwnerId,
			p.TransactionId,
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

func (repo *PromoCodeJointRepo) UpdatePromoCodeJoinTxByTransactionId(ctx context.Context, p PromoCodeJoint) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	if ok {
		_, err = dTx.Exec(sQLPromoCodeJointUpdateByTransactionId, p.Status, p.TransactionId)

	} else {
		_, err = repo.db.Exec(sQLPromoCodeJointUpdateByTransactionId, p.Status, p.TransactionId)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *PromoCodeJointRepo) GetAllActivePromoCodeJointByOwnerAndPromoCode(ctx context.Context, ownerId int64, promoCode string) ([]*PromoCodeJoint, error) {
	rows, err := repo.db.Query(sQLPromoCodeJointGetAllActiveByOwnerAndPromoCode, ownerId, promoCode)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pJoint := make([]*PromoCodeJoint, 16)
	for rows.Next() {
		p, err := scanRowIntoPromoCodeJoint(rows)
		if err != nil {
			return nil, err
		}
		pJoint = append(pJoint, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pJoint, nil
}

func scanRowIntoPromoCodeJoint(rows *sql.Rows) (*PromoCodeJoint, error) {
	p := new(PromoCodeJoint)
	err := rows.Scan(
		&p.ID,
		&p.Status,
		&p.PromoCode,
		&p.OwnerId,
		&p.TransactionId,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}
