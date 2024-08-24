package approval

import "time"

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
		LIMIT ? OFFSET ?
	`
	sQLApprovalQueryAllByTypeAndStatusWithPagination = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.decided_by, a.decided_at,
			a.create_at, a.update_at
		FROM approval as a
		WHERE a.type = ? AND a.status = ? 
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
