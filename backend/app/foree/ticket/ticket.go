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
	ID          int64        `json:"id"`
	Type        string       `json:"type"`
	Priority    string       `json:"priority"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Status      TicketStatus `json:"status"`
	ConcludedBy int64        `json:"concludedBy"`
	ConcludedAt time.Time    `json:"concludedAt"`
	CreateAt    time.Time    `json:"createAt"`
	UpdateAt    time.Time    `json:"updateAt"`
}

type TicketRepo struct {
	db *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
	return &TicketRepo{db: db}
}
