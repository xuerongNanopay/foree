package approval

import (
	"database/sql"
	"time"
)

const (
	sQLApprovalInsert = `
		INSERT INTO approval
		(
			type, status, associated_entity_name,
			associated_entity_id
		) VALUES (?,?,?,?)
	`
	sQLApprovalUpdate = `
		UPDATE approval SET
			status = ?, decided_by = ?, decided_at = ?
		WHERE id = ?
	`
	sQLApprovalGetUniqueById = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.decided_by, a.decided_at,
			a.create_at, a.update_at
		FROM approval as a
		WHERE a.id = ?
	`
	sQLApprovalQueryAllByTypeWithPagination = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.decided_by, a.decided_at,
			a.create_at, a.update_at
		FROM approval as a
		WHERE a.type = ?
		ORDER BY a.create_at DESC
		LIMIT ? OFFSET ?
	`
	sQLApprovalQueryAllByTypeAndStatusWithPagination = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.decided_by, a.decided_at,
			a.create_at, a.update_at
		FROM approval as a
		WHERE a.type = ? AND a.status = ? 
		ORDER BY a.create_at DESC
		LIMIT ? OFFSET ?
	`
)

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "PENDING"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "Rejected"
)

type Approval struct {
	ID                   int64          `json:"id"`
	Type                 string         `json:"type"`
	Status               ApprovalStatus `json:"status"`
	AssociatedEntityName string         `json:"associatedEntityName"`
	AssocitateEntityId   int64          `json:"associtateEntityId"`
	DecidedBy            int64          `json:"decidedBy"`
	DecidedAt            time.Time      `json:"decidedAt"`
	CreateAt             time.Time      `json:"createAt"`
	UpdateAt             time.Time      `json:"updateAt"`
}

type ApprovalRepo struct {
	db *sql.DB
}

func NewApprovalRepo(db *sql.DB) *ApprovalRepo {
	return &ApprovalRepo{db: db}
}

func (repo *ApprovalRepo) InsertApproval(approval Approval) (int64, error) {
	result, err := repo.db.Exec(
		sQLApprovalInsert,
		approval.Type,
		approval.Status,
		approval.AssociatedEntityName,
		approval.AssocitateEntityId,
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

func (repo *ApprovalRepo) UpdateApprovalById(approval Approval) error {
	_, err := repo.db.Exec(
		sQLApprovalUpdate,
		approval.Status,
		approval.DecidedBy,
		approval.DecidedAt,
		approval.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoApproval(rows *sql.Rows) (*Approval, error) {
	i := new(Approval)
	err := rows.Scan(
		&i.ID,
		&i.Type,
		&i.Status,
		&i.AssociatedEntityName,
		&i.AssocitateEntityId,
		&i.DecidedBy,
		&i.DecidedAt,
		&i.CreateAt,
		&i.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}
