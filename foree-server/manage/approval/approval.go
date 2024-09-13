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
			status = ?, concluded_by = ?, concluded_at = ?
		WHERE id = ?
	`
	sQLApprovalGetUniqueById = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.concluded_by, a.concluded_at,
			a.created_at, a.updated_at
		FROM approval as a
		WHERE a.id = ?
	`
	sQLApprovalQueryAllByTypeWithPagination = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.concluded_by, a.concluded_at,
			a.created_at, a.updated_at
		FROM approval as a
		WHERE a.type = ?
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`
	sQLApprovalQueryAllByTypeAndStatusWithPagination = `
		SELECT
			a.id, a.type, a.status, a.associated_entity_name,
			a.associated_entity_id, a.concluded_by, a.concluded_at,
			a.created_at, a.updated_at
		FROM approval as a
		WHERE a.type = ? AND a.status = ? 
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`
)

type ApprovalStatus string

const (
	ApprovalStatusOpen     ApprovalStatus = "OPEN"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

type Approval struct {
	ID                   int64          `json:"id"`
	Type                 string         `json:"type"`
	Status               ApprovalStatus `json:"status"`
	AssociatedEntityName string         `json:"associatedEntityName"`
	AssocitateEntityId   int64          `json:"associtateEntityId"`
	ConcludedBy          int64          `json:"concludedBy"`
	ConcludedAt          time.Time      `json:"concludedAt"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
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
		approval.ConcludedBy,
		approval.ConcludedAt,
		approval.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ApprovalRepo) GetUniqueApprovalById(id int64) (*Approval, error) {
	rows, err := repo.db.Query(sQLApprovalGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *Approval

	for rows.Next() {
		u, err = scanRowIntoApproval(rows)
		if err != nil {
			return nil, err
		}
	}

	if u == nil || u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func (repo *ApprovalRepo) QueryAllApprovalByTypeWithPagination(t string, limit, offset int) ([]*Approval, error) {
	rows, err := repo.db.Query(sQLApprovalQueryAllByTypeWithPagination, t, limit, offset)

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

func (repo *ApprovalRepo) QueryAllApprovalByTypeAndStatusWithPagination(t, s string, limit, offset int) ([]*Approval, error) {
	rows, err := repo.db.Query(sQLApprovalQueryAllByTypeAndStatusWithPagination, t, s, limit, offset)

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

func scanRowIntoApproval(rows *sql.Rows) (*Approval, error) {
	i := new(Approval)
	err := rows.Scan(
		&i.ID,
		&i.Type,
		&i.Status,
		&i.AssociatedEntityName,
		&i.AssocitateEntityId,
		&i.ConcludedBy,
		&i.ConcludedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}
