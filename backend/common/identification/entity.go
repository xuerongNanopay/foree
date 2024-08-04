package identification

import "time"

type IdentificationType string

const (
	IDTypePassport      = "PASSPORT"
	IDTypeDriverLicense = "DRIVER_LICENSE"
	IDTypeProvincalId   = "PROVINCIAL_ID"
	IDTypeNationId      = "NATIONAL_ID"
)

type Identification struct {
	ID           uint64             `json:"id"`
	Type         IdentificationType `json:"type"`
	Value        string             `json:"value"`
	AdditionInfo string             `json:"additionInfo"`
	CreateAt     time.Timer         `json:"-"`
	UpdateAt     time.Timer         `json:"-"`
	OwnerId      int64              `json:"ownerId"`
}
