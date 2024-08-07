package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLInteracCITxInsert = `
		INSERT INTO interact_ci_tx
		(
			status, src_interac_acc_id, desc_interac_acc_
		) VALUES()
	`
	sQLInteracCITxGetUniqueByParentTxId = `
	
	`
	sQLInteracCITxGetUniqueById = `
	
	`
	sQLInteracCITxUpdateById = `
	
	`
)

type ScotiaInteracCITransaction struct {
	ID                  int64
	Status              TxStatus
	ScotialId           string
	Url                 string
	SrcInteracAccId     int64
	SrcInteracAcc       *ScotiaInteracCITransaction
	DescInteracAccId    int64
	DescInteracAcc      *ScotiaInteracCITransaction
	Amt                 types.AmountData
	ParentTransactionId int64
	OwnerId             int64
	CreateAt            time.Time `json:"createAt"`
	UpdateAt            time.Time `json:"updateAt"`
}

func NewInteracCIRepo(db *sql.DB) *InteracCIRepo {
	return &InteracCIRepo{db: db}
}

type InteracCIRepo struct {
	db *sql.DB
}
