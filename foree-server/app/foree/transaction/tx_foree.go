package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
)

const (
	sQLForeeTxInsert = `
		INSERT INTO foree_tx
		(
			type, rate, cin_acc_id, cout_acc_id, limit_reference,
			src_amount, src_currency, dest_amount, dest_currency,
			total_fee_amount, total_fee_currency, total_reward_amount, total_reward_currency,
			total_amount, total_currency, stage,
			transaction_purpose, conclusion, owner_id
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLForeeTxUpdateById = `
	    UPDATE foree_tx SET 
			stage = ?, conclusion = ?
        WHERE id = ?
	`
	sQLForeeTxGetById = `
	    SELECT 
            t.id, t.type, t.rate,
			t.cin_acc_id, t.cout_acc_id, t.limit_reference,
            t.src_amount, t.src_currency, 
            t.dest_amount, t.dest_currency,
			t.total_fee_amount, t.total_fee_currency, 
            t.total_reward_amount, t.total_reward_currency, 
            t.total_amount, t.total_currency,
			t.stage, t.transaction_purpose, t.conclusion,
            t.owner_id, t.created_at, t.updated_at
        FROM foree_tx t
        where t.id = ?
	`
	sQLForeeTxGetForUpdateById = `
		SELECT 
			t.id, t.type, t.rate,
			t.cin_acc_id, t.cout_acc_id, t.limit_reference,
			t.src_amount, t.src_currency, 
			t.dest_amount, t.dest_currency,
			t.total_fee_amount, t.total_fee_currency, 
			t.total_reward_amount, t.total_reward_currency, 
			t.total_amount, t.total_currency,
			t.stage, t.transaction_purpose, t.conclusion,
			t.owner_id, t.created_at, t.updated_at
		FROM foree_tx t
		where t.id = ?
		FOR UPDATE
	`
	//TODO: support get alls?
)

type TxStatus string

const (
	TxStatusInitial    TxStatus = "INITIAL"
	TxStatusProcessing TxStatus = "PROCESSING"
	TxStatusSuspend    TxStatus = "SUSPEND"
	TxStatusSent       TxStatus = "SENT"
	TxStatusRejected   TxStatus = "REJECTED"
	TxStatusCancelled  TxStatus = "CANCELLED"
	TxStatusCompleted  TxStatus = "COMPLETED"
	TxStatusClosed     TxStatus = "CLOSED"
)

type TxType string

const (
	TxTypeInteracToNBP TxType = "INTERAC-NBP"
)

type TxStage string

const (
	TxStageBegin     TxStage = "TX_BEGIN"
	TxStageInteracCI TxStage = "CASHIN_INTERAC"
	TxStageIDM       TxStage = "COMPLIANCE_IDM"
	TxStageNBPCO     TxStage = "CASHOUT_NBP"
	TxStageRefunding TxStage = "TX_REFUNDING"
	TxStageCancel    TxStage = "TX_CANCEL"
	TxStageSuccess   TxStage = "TX_SUCCESS"
)

// Only Support INITIAL, PROCESSING, CANCEL, COMPLETE.
type ForeeTx struct {
	ID                 int64            `json:"id,omitempty"`
	Type               TxType           `json:"type,omitempty"`
	Rate               types.Amount     `json:"Rate,omitempty"`
	CinAccId           int64            `json:"cinAccId,omitempty"`
	CoutAccId          int64            `json:"coutAccId,omitempty"`
	LimitReference     string           `json:"limitReference,omitempty"`
	SrcAmt             types.AmountData `json:"srcAmt,omitempty"`
	DestAmt            types.AmountData `json:"destAmt,omitempty"`
	TotalFeeAmt        types.AmountData `json:"totalFeeAmt,omitempty"`
	TotalRewardAmt     types.AmountData `json:"totalRewardAmt,omitempty"`
	TotalAmt           types.AmountData `json:"totalAmt,omitempty"`
	Stage              TxStage          `json:"curStage,omitempty"`
	TransactionPurpose string           `json:"transactionPurpose,omitempty"`
	Conclusion         string           `json:"conclusion,omitempty"`
	OwnerId            int64            `json:"ownerId,omitempty"`
	CreatedAt          *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time       `json:"updatedAt,omitempty"`

	Ip            string                  `json:"ip,omitempty"`
	UserAgent     string                  `json:"userAgent,omitempty"`
	Owner         *auth.User              `json:"ower,omitempty"`
	InteracAcc    *account.InteracAccount `json:"interacAcc,omitempty"`
	ContactAcc    *account.ContactAccount `json:"contactAcc,omitempty"`
	FeeJointIds   []int64                 `json:"feeJointIds,omitempty"`
	FeeJoints     []*FeeJoint             `json:"feesJoints,omitempty"`
	RewardIds     []int64                 `json:"rewardIds,omitempty"`
	Rewards       []*promotion.Reward     `json:"rewards,omitempty"`
	CI            *InteracCITx            `json:"ci,omitempty"`
	IDM           *IDMTx                  `json:"idm,omitempty"`
	COUT          *NBPCOTx                `json:"cout,omitempty"`
	Summary       *TxSummary              `json:"summary,omitempty"`
	ForeeRefundTx *ForeeRefundTx          `json:"refundTx,omitempty"`
	History       []*TxHistory            `json:"history,omitempty"`
}

func NewForeeTxRepo(db *sql.DB) *ForeeTxRepo {
	return &ForeeTxRepo{db: db}
}

type ForeeTxRepo struct {
	db *sql.DB
}

func (repo *ForeeTxRepo) InsertForeeTx(ctx context.Context, tx ForeeTx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLForeeTxInsert,
			tx.Type,
			tx.Rate,
			tx.CinAccId,
			tx.CoutAccId,
			tx.LimitReference,
			tx.SrcAmt.Amount,
			tx.SrcAmt.Currency,
			tx.DestAmt.Amount,
			tx.DestAmt.Currency,
			tx.TotalFeeAmt.Amount,
			tx.TotalFeeAmt.Currency,
			tx.TotalRewardAmt.Amount,
			tx.TotalRewardAmt.Currency,
			tx.TotalAmt.Amount,
			tx.TotalAmt.Currency,
			tx.Stage,
			tx.TransactionPurpose,
			tx.Conclusion,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLForeeTxInsert,
			tx.Type,
			tx.Rate,
			tx.CinAccId,
			tx.CoutAccId,
			tx.LimitReference,
			tx.SrcAmt.Amount,
			tx.SrcAmt.Currency,
			tx.DestAmt.Amount,
			tx.DestAmt.Currency,
			tx.TotalFeeAmt.Amount,
			tx.TotalFeeAmt.Currency,
			tx.TotalRewardAmt.Amount,
			tx.TotalRewardAmt.Currency,
			tx.TotalAmt.Amount,
			tx.TotalAmt.Currency,
			tx.Stage,
			tx.TransactionPurpose,
			tx.Conclusion,
			tx.OwnerId,
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

func (repo *ForeeTxRepo) UpdateForeeTxById(ctx context.Context, tx ForeeTx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	if ok {
		_, err = dTx.Exec(sQLForeeTxUpdateById, tx.Stage, tx.Conclusion, tx.ID)
	} else {
		_, err = repo.db.Exec(sQLForeeTxUpdateById, tx.Stage, tx.Conclusion, tx.ID)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *ForeeTxRepo) GetUniqueForeeTxById(ctx context.Context, id int64) (*ForeeTx, error) {
	rows, err := repo.db.Query(sQLForeeTxGetById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ForeeTx

	for rows.Next() {
		f, err = scanRowIntoForeeTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ForeeTxRepo) GetUniqueForeeTxForUpdateById(ctx context.Context, id int64) (*ForeeTx, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var rows *sql.Rows

	if ok {
		rows, err = dTx.Query(sQLForeeTxGetForUpdateById, id)

	} else {
		rows, err = repo.db.Query(sQLForeeTxGetForUpdateById, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ForeeTx

	for rows.Next() {
		f, err = scanRowIntoForeeTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoForeeTx(rows *sql.Rows) (*ForeeTx, error) {
	tx := new(ForeeTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Type,
		&tx.Rate,
		&tx.CinAccId,
		&tx.CoutAccId,
		&tx.LimitReference,
		&tx.SrcAmt.Amount,
		&tx.SrcAmt.Currency,
		&tx.DestAmt.Amount,
		&tx.DestAmt.Currency,
		&tx.TotalFeeAmt.Amount,
		&tx.TotalFeeAmt.Currency,
		&tx.TotalRewardAmt.Amount,
		&tx.TotalRewardAmt.Currency,
		&tx.TotalAmt.Amount,
		&tx.TotalAmt.Currency,
		&tx.Stage,
		&tx.TransactionPurpose,
		&tx.Conclusion,
		&tx.OwnerId,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
