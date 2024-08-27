package ticket

import (
	"database/sql"
	"time"
)

type TicketStatus string

const (
	TicketStatusOpen     TicketStatus = "OPEN"
	TicketStatusResolved TicketStatus = "RESOLVED"
	TicketStatusClosed   TicketStatus = "CLOSED"
)

type Ticket struct {
	ID                   int64        `json:"id"`
	Type                 string       `json:"type"`
	Status               TicketStatus `json:"status"`
	Priority             string       `json:"priority"`
	Summary              string       `json:"summary"`
	Description          string       `json:"description"`
	AssociatedEntityName string       `json:"associatedEntityName"`
	AssocitateEntityId   int64        `json:"associtateEntityId"`
	ConcludedBy          int64        `json:"concludedBy"`
	ConcludedAt          time.Time    `json:"concludedAt"`
	CreateAt             time.Time    `json:"createAt"`
	UpdateAt             time.Time    `json:"updateAt"`
}

type TicketRepo struct {
	db *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
	return &TicketRepo{db: db}
}
