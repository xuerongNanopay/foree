package approval

import "time"

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
