package idm_compliance

import "time"

type IDMCompliance struct {
	ID             int64     `json:"id"`
	ForeeTxId      int64     `json:"foreeTxId"`
	IDMRespStatus  string    `json:"idmRespStatus"`
	IDMResHttpCode string    `json:"idmRespHttpCode"`
	IDMRawRequest  string    `json:"idmRawRequest"`
	IDMRawResponse string    `json:"idmRawResponse"`
	CreateAt       time.Time `json:"createAt"`
	UpdateAt       time.Time `json:"updateAt"`
}
