package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
)

type TxSummaryStatus string

const (
	TxSummaryStatusInitial      TxSummaryStatus = "Initial"
	TxSummaryStatusAwaitPayment TxSummaryStatus = "Await Payment"
	TxSummaryStatusInProgress   TxSummaryStatus = "In Progress"
	TxSummaryStatusCompleted    TxSummaryStatus = "Completed"
	TxSummaryStatusCancelled    TxSummaryStatus = "Cancelled"
	TxSummaryStatusPickup       TxSummaryStatus = "Ready To Pickup"
	TxSummaryStatusRefunding    TxSummaryStatus = "Refunding"
	TxSummaryStatusRefunded     TxSummaryStatus = "Refunded"
)

const (
	sQLTxSummaryInsert = `
        INSERT INTO tx_summary
        (
            summary, type, status, rate, 
			src_acc_id, dest_acc_id,
            src_acc_summary, src_amount, src_currency, 
            dest_acc_summary, dest_amount, dest_currency,
            total_amount, total_currency,
            fee_amount, fee_currency, 
            reward_amount, reward_currency, 
            nbp_reference, is_cancel_allowed, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
    `
	sQLTxSummaryUpdateById = `
        UPDATE tx_summary SET 
            status = ?, is_cancel_allowed = ?, payment_url = ?
        WHERE id = ?
    `
	sQLTxSummaryGetUniqueByOwnerAndId = `
        SELECT 
            t.id, t.summary, t.type, t.status, t.rate,
			t.payment_url, t.src_acc_id, t.dest_acc_id,
            t.src_acc_summary, t.src_amount, t.src_currency, 
            t.dest_acc_summary, t.dest_amount, t.dest_currency,
            t.total_amount, t.total_currency,
            t.fee_amount, t.fee_currency, 
            t.reward_amount, t.reward_currency, 
            t.nbp_reference, t.is_cancel_allowed, t.parent_tx_id, t.owner_id, 
            t.created_at, t.updated_at
        FROM tx_summary t
        where t.owner_id = ? and t.id = ?
    `
	sQLTxSummaryGetUniqueByParentTxId = `
        SELECT 
            t.id, t.summary, t.type, t.status, t.rate,
			t.payment_url, t.src_acc_id, t.dest_acc_id,
            t.src_acc_summary, t.src_amount, t.src_currency, 
            t.dest_acc_summary, t.dest_amount, t.dest_currency,
            t.total_amount, t.total_currency,
            t.fee_amount, t.fee_currency, 
            t.reward_amount, t.reward_currency, 
            t.nbp_reference, t.is_cancel_allowed, t.parent_tx_id, t.owner_id, 
            t.created_at, t.updated_at
        FROM tx_summary t
        where t.parent_tx_id = ?
    `
	sQLTxSummaryGetAllByOwnerIdWithPagination = `
	    SELECT
	        t.id, t.summary, t.type, t.status, t.rate,
			t.payment_url, t.src_acc_id, t.dest_acc_id,
	        t.src_acc_summary, t.src_amount, t.src_currency,
	        t.dest_acc_summary, t.dest_amount, t.dest_currency,
	        t.total_amount, t.total_currency,
	        t.fee_amount, t.fee_currency,
	        t.reward_amount, t.reward_currency,
	        t.nbp_reference, t.is_cancel_allowed, t.parent_tx_id, t.owner_id,
	        t.created_at, t.updated_at
	    FROM tx_summary t
	    where t.owner_id = ?
	    ORDER BY t.created_at DESC
	    LIMIT ? OFFSET ?
	`
	sQLTxSummaryQueryByOwnerIdAndStatusWithPagination = `
	    SELECT
	        t.id, t.summary, t.type, t.status, t.rate,
			t.payment_url, t.src_acc_id, t.dest_acc_id,
	        t.src_acc_summary, t.src_amount, t.src_currency,
	        t.dest_acc_summary, t.dest_amount, t.dest_currency,
	        t.total_amount, t.total_currency,
	        t.fee_amount, t.fee_currency,
	        t.reward_amount, t.reward_currency,
	        t.nbp_reference, t.is_cancel_allowed, t.parent_tx_id, t.owner_id,
	        t.created_at, t.updated_at
	    FROM tx_summary t
	    where t.owner_id = ? AND t.status = ?
	    ORDER BY t.created_at DESC
	    LIMIT ? OFFSET ?
	`
)

type TxSummary struct {
	ID              int64           `json:"id"`
	Summary         string          `json:"sumary"`
	Type            string          `json:"type"`
	Status          TxSummaryStatus `json:"status"`
	Rate            string          `json:"rate"`
	PaymentUrl      string          `json:"paymentUrl"`
	SrcAccId        int64           `json:"srcAccId"`
	DestAccId       int64           `json:"destAccId"`
	SrcAccSummary   string          `json:"srcAccSummary"`
	SrcAmount       types.Amount    `json:"srcAmount"`
	SrcCurrency     string          `json:"srcCurrency"`
	DestAccSummary  string          `json:"destAccSummary"`
	DestAmount      types.Amount    `json:"destAmount"`
	DestCurrency    string          `json:"destCurrency"`
	TotalAmount     types.Amount    `json:"totalAmount"`
	TotalCurrency   string          `json:"totalCurrency"`
	FeeAmount       types.Amount    `json:"feeAmount"`
	FeeCurrency     string          `json:"feeCurrency"`
	RewardAmount    types.Amount    `json:"rewardAmount"`
	RewardCurrency  string          `json:"rewardCurrency"`
	NBPReference    string          `json:"nbpReference"`
	IsCancelAllowed bool            `json:"isCancelAllowed"`
	ParentTxId      int64           `json:"parentTxd"`
	OwnerId         int64           `json:"owerId"`
	CreatedAt       *time.Time      `json:"createdAt"`
	UpdatedAt       *time.Time      `json:"updatedAt"`

	SrcAccount  *account.InteracAccount
	DestAccount *account.ContactAccount
}

func NewTxSummaryRepo(db *sql.DB) *TxSummaryRepo {
	return &TxSummaryRepo{db: db}
}

type TxSummaryRepo struct {
	db *sql.DB
}

func (repo *TxSummaryRepo) InsertTxSummary(ctx context.Context, tx TxSummary) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result
	if ok {
		result, err = dTx.Exec(
			sQLTxSummaryInsert,
			tx.Summary,
			tx.Type,
			tx.Status,
			tx.Rate,
			tx.SrcAccId,
			tx.DestAccId,
			tx.SrcAccSummary,
			tx.SrcAmount,
			tx.SrcCurrency,
			tx.DestAccSummary,
			tx.DestAmount,
			tx.DestCurrency,
			tx.TotalAmount,
			tx.TotalCurrency,
			tx.FeeAmount,
			tx.FeeCurrency,
			tx.RewardAmount,
			tx.RewardCurrency,
			tx.NBPReference,
			tx.IsCancelAllowed,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLTxSummaryInsert,
			tx.Summary,
			tx.Type,
			tx.Status,
			tx.Rate,
			tx.SrcAccId,
			tx.DestAccId,
			tx.SrcAccSummary,
			tx.SrcAmount,
			tx.SrcCurrency,
			tx.DestAccSummary,
			tx.DestAmount,
			tx.DestCurrency,
			tx.TotalAmount,
			tx.TotalCurrency,
			tx.FeeAmount,
			tx.FeeCurrency,
			tx.RewardAmount,
			tx.RewardCurrency,
			tx.NBPReference,
			tx.IsCancelAllowed,
			tx.ParentTxId,
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

func (repo *TxSummaryRepo) UpdateTxSummaryById(ctx context.Context, tx TxSummary) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	if ok {
		_, err = dTx.Exec(sQLTxSummaryUpdateById, tx.Status, tx.IsCancelAllowed, tx.PaymentUrl, tx.ID)
	} else {
		_, err = repo.db.Exec(sQLTxSummaryUpdateById, tx.Status, tx.IsCancelAllowed, tx.PaymentUrl, tx.ID)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *TxSummaryRepo) GetUniqueTxSummaryByOwnerAndId(ctx context.Context, userId, id int64) (*TxSummary, error) {
	rows, err := repo.db.Query(sQLTxSummaryGetUniqueByOwnerAndId, userId, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *TxSummary

	for rows.Next() {
		f, err = scanRowIntoTxSummary(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *TxSummaryRepo) GetUniqueTxSummaryByParentTxId(ctx context.Context, parentTxId int64) (*TxSummary, error) {
	rows, err := repo.db.Query(sQLTxSummaryGetUniqueByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *TxSummary

	for rows.Next() {
		f, err = scanRowIntoTxSummary(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *TxSummaryRepo) GetAllTxSummaryByOwnerIdWithPagination(ctx context.Context, ownerId int64, limit, offset int) ([]*TxSummary, error) {
	rows, err := repo.db.Query(sQLTxSummaryGetAllByOwnerIdWithPagination, ownerId, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*TxSummary, 0)
	for rows.Next() {
		p, err := scanRowIntoTxSummary(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (repo *TxSummaryRepo) QueryTxSummaryByOwnerIdAndStatusWithPagination(ctx context.Context, ownerId int64, status string, limit, offset int) ([]*TxSummary, error) {
	rows, err := repo.db.Query(sQLTxSummaryQueryByOwnerIdAndStatusWithPagination, ownerId, status, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*TxSummary, 0)
	for rows.Next() {
		p, err := scanRowIntoTxSummary(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func scanRowIntoTxSummary(rows *sql.Rows) (*TxSummary, error) {
	tx := new(TxSummary)
	err := rows.Scan(
		&tx.ID,
		&tx.Summary,
		&tx.Type,
		&tx.Status,
		&tx.Rate,
		&tx.PaymentUrl,
		&tx.SrcAccId,
		&tx.DestAccId,
		&tx.SrcAccSummary,
		&tx.SrcAmount,
		&tx.SrcCurrency,
		&tx.DestAccSummary,
		&tx.DestAmount,
		&tx.DestCurrency,
		&tx.TotalAmount,
		&tx.TotalCurrency,
		&tx.FeeAmount,
		&tx.FeeCurrency,
		&tx.RewardAmount,
		&tx.RewardCurrency,
		&tx.NBPReference,
		&tx.IsCancelAllowed,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
