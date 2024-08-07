package transaction

import (
	"time"
)

type IDMTransaction struct {
	ID         int64
	Status     TxStatus
	Ip         string    `json:"ip"`
	UserAgent  string    `json:"userAgent"`
	IDMResult  string    `json:"idmResult"`
	ParentTxId int64     `json:"parentTxId"`
	OwnerId    int64     `json:"ownerId"`
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}

// Large object.
type IDMCompliance struct {
	ID                int64     `json:"id"`
	IDMTxId           int64     `json:"idmTxId"`
	IDMHttpStatusCode int       `json:"idmHttpStatusCode"`
	IDMResult         string    `json:"idmResult"`
	RequestJson       string    `json:"requestJson"`
	ResponseJson      string    `json:"responseJson"`
	CreateAt          time.Time `json:"createAt"`
	UpdateAt          time.Time `json:"updateAt"`
}
