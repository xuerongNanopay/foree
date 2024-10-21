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
	sQLApprovalUpdateById = `
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
	sQLApprovalGetAllWithPagination = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
	    ORDER BY a.created_at DESC
	    LIMIT ? OFFSET ?
	`
	sQLApprovalGetAllByTypeWithPagination = `
		SELECT
			a.id, a.type, a.description, a.status, a.ref_entity, a.ref_id,
			a.approved_by, a.rejected_by, a.reject_reason, a.approved_at,
			a.rejected_at, a.owner_id, a.created_at, a.updated_at
		FROM approval a
		WHERE a.type = ?
	    ORDER BY a.created_at DESC
	    LIMIT ? OFFSET ?
	`
	sQLApprovalGetAllByTypeAndStatusWithPagination = `
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

func (repo *ApprovalRepo) UpdateApprovalById(ctx context.Context, a Approval) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	if ok {
		_, err = dTx.ExecContext(
			ctx,
			sQLApprovalUpdateById,
			a.Status,
			a.ApprovedBy,
			a.RejectedBy,
			a.RejectReason,
			a.ApprovedAt,
			a.RejectedAt,
			a.OwnerId,
		)
	} else {
		_, err = repo.db.ExecContext(
			ctx,
			sQLApprovalUpdateById,
			a.Status,
			a.ApprovedBy,
			a.RejectedBy,
			a.RejectReason,
			a.ApprovedAt,
			a.RejectedAt,
			a.OwnerId,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *ApprovalRepo) GetUniqueApprovalyById(ctx context.Context, id int64) (*Approval, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var rows *sql.Rows

	if ok {
		rows, err = dTx.Query(sQLApprovalGetUniqueApprovalById, id)
	} else {
		rows, err = repo.db.Query(sQLApprovalGetUniqueApprovalById, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Approval

	for rows.Next() {
		f, err = scanRowIntoApproval(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ApprovalRepo) getAll(ctx context.Context, query string, args ...any) ([]*Approval, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(query, args)
	} else {
		rows, err = repo.db.Query(query, args)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvals := make([]*Approval, 0)
	for rows.Next() {
		p, err := scanRowIntoApproval(rows)
		if err != nil {
			return nil, err
		}
		approvals = append(approvals, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return approvals, nil
}

func (repo *ApprovalRepo) GetAllApprovalWithPagination(ctx context.Context, limit, offset int) ([]*Approval, error) {
	return repo.getAll(ctx, sQLApprovalGetAllWithPagination, limit, offset)
}

func (repo *ApprovalRepo) GetAllApprovalByTypeWithPagination(ctx context.Context, approvalType string, limit, offset int) ([]*Approval, error) {
	return repo.getAll(ctx, sQLApprovalGetAllByTypeWithPagination, approvalType, limit, offset)
}

func (repo *ApprovalRepo) GetAllApprovalByTypeAndStatusWithPagination(ctx context.Context, approvalType, status string, limit, offset int) ([]*Approval, error) {
	return repo.getAll(ctx, sQLApprovalGetAllByTypeAndStatusWithPagination, approvalType, status, limit, offset)
}

func scanRowIntoApproval(rows *sql.Rows) (*Approval, error) {
	a := new(Approval)
	err := rows.Scan(
		&a.ID,
		&a.Type,
		&a.Description,
		&a.Status,
		&a.RefEntity,
		&a.RefId,
		&a.ApprovedBy,
		&a.RejectedBy,
		&a.RejectReason,
		&a.ApprovedAt,
		&a.RejectedAt,
		&a.OwnerId,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}
