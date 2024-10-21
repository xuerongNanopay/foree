package approval

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

type ApproveStatus string

const (
	ApprovalStatusPending  ApproveStatus = "PENDING"
	ApprovalStatusRejected ApproveStatus = "REJECTED"
	ApprovalStatusApproved ApproveStatus = "APPROVED"
)

const (
	sQLApprovalInsert = `
		INSERT INTO approval(
			type, description, status, ref_entity, ref_id, owner_id
		) VALUES(?,?,?,?,?,?)
	`
	sQLApprovalUpdate = `
		UPDATE approval SET
			status = ?, approved_by = ?, rejected_by = ?, reject_reason = ?,
			approved_at = ?, rejected_at = ?
		WHERE id = ?
	`
	sQLApprovalGetUniqueApprovalById = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
		WHERE a.id = ?
	`
	sQLApprovalGetAllApprovalWithPagination = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
	    ORDER BY a.created_at DESC
	    LIMIT ? OFFSET ?
	`
	sQLApprovalGetAllApprovalByTypeWithPagination = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
		WHERE a.type = ?
	    ORDER BY a.created_at DESC
	    LIMIT ? OFFSET ?
	`
	sQLApprovalGetAllApprovalByTypeAndStatusWithPagination = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
		WHERE a.type = ? AND a.status = ?
	    ORDER BY a.created_at DESC
	    LIMIT ? OFFSET ?
	`
)

type Approval struct {
	ID           int64         `json:"id"`
	Type         string        `json:"type"`
	Description  string        `json:"description"`
	Status       ApproveStatus `json:"status"`
	RefEntity    string        `json:"refEntity"`
	RefId        int64         `json:"refId"`
	ApprovedBy   int64         `json:"approvedBy"`
	RejectedBy   int64         `json:"rejectedBy"`
	RejectReason string        `json:"rejectReason"`
	ApprovedAt   *time.Time    `json:"approvedAt"`
	RejectedAt   *time.Time    `json:"rejectedAt"`
	OwnerId      int64         `json:"ownerId"`
	CreatedAt    *time.Time    `json:"createdAt"`
	UpdatedAt    *time.Time    `json:"updatedAt"`
}

func NewApprovalRepo(db *sql.DB) *ApprovalRepo {
	return &ApprovalRepo{db: db}
}

type ApprovalRepo struct {
	db *sql.DB
}

func (repo *ApprovalRepo) InsertApproval(ctx context.Context, a Approval) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)
	var err error
	var result sql.Result

	if ok {
		result, err = dTx.ExecContext(
			ctx,
			sQLApprovalInsert,
			a.Type,
			a.Description,
			a.Status,
			a.RefEntity,
			a.RefId,
			a.OwnerId,
		)
	} else {
		result, err = repo.db.ExecContext(
			ctx,
			sQLApprovalInsert,
			a.Type,
			a.Description,
			a.Status,
			a.RefEntity,
			a.RefId,
			a.OwnerId,
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
