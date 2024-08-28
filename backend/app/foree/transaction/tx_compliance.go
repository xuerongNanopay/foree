package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLIDMTxInsert = `
        INSERT INTO idm_tx
        (
            status, ip, user_agent, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?)
    `
	sQLIDMTxUpdateById = `
        UPDATE idm_tx SET 
            status = ?, idm_reference = ?, t.idm_result = ?
        WHERE id = ?
    `
	sQLIDMTxGetUniqueById = `
        SELECT 
            t.id, t.status, t.ip, t.user_agent, t.idm_reference, t.idm_result
            t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
        FROM idm_tx t
        where t.id = ?
    `
	sQLIDMTxGetUniqueByParentTxId = `
        SELECT 
            t.id, t.status, t.ip, t.user_agent, t.idm_reference, t.idm_result
            t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
        FROM idm_tx t
        where t.parent_tx_id = ?
    `
	sQLIDMComplianceInsert = `
        INSERT INTO idm_compliance
        (
            idm_tx_id, idm_http_status_code, idm_result, request_json, response_json
        ) VALUES(?,?,?,?,?)
    `
	sQLIDMComplianceGetUniqueById = `
        SELECT 
            c.id, c.idm_tx_id, c.idm_http_status_code, c.idm_result, 
            c.request_json, c.response_json,
            c.created_at, c.updated_at
        FROM idm_tx c
        where c.id = ?
    `
	sQLIDMComplianceGetUniqueByIDMTxId = `
        SELECT 
            c.id, c.idm_tx_id, c.idm_http_status_code, c.idm_result, 
            c.request_json, c.response_json,
            c.created_at, c.updated_at
        FROM idm_tx c
        where c.idm_tx_id = ?
    `
)

type IDMTx struct {
	ID           int64
	Status       TxStatus
	Ip           string    `json:"ip"`
	UserAgent    string    `json:"userAgent"`
	IDMReference string    `json:"idmReference"`
	IDMResult    string    `json:"idmResult"`
	ParentTxId   int64     `json:"parentTxId"`
	OwnerId      int64     `json:"ownerId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdateAt     time.Time `json:"updatedAt"`
}

// Large object.
type IDMCompliance struct {
	ID                int64     `json:"id"`
	IDMTxId           int64     `json:"idmTxId"`
	IDMHttpStatusCode int       `json:"idmHttpStatusCode"`
	IDMResult         string    `json:"idmResult"`
	RequestJson       string    `json:"requestJson"`
	ResponseJson      string    `json:"responseJson"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdateAt          time.Time `json:"updatedAt"`
}

func NewIdmTxRepo(db *sql.DB) *IdmTxRepo {
	return &IdmTxRepo{db: db}
}

type IdmTxRepo struct {
	db *sql.DB
}

func (repo *IdmTxRepo) InsertIDMTx(ctx context.Context, tx IDMTx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLIDMTxInsert,
			tx.Status,
			tx.Ip,
			tx.UserAgent,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLIDMTxInsert,
			tx.Status,
			tx.Ip,
			tx.UserAgent,
			tx.ParentTxId,
			tx.OwnerId,
		)
	}

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *IdmTxRepo) UpdateIDMTxById(ctx context.Context, tx IDMTx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLIDMTxUpdateById, tx.Status, tx.IDMReference, tx.IDMResult, tx.ID)

	} else {
		_, err = repo.db.Exec(sQLIDMTxUpdateById, tx.Status, tx.IDMReference, tx.IDMResult, tx.ID)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *IdmTxRepo) GetUniqueIDMTxById(ctx context.Context, id int64) (*IDMTx, error) {
	rows, err := repo.db.Query(sQLIDMTxGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *IDMTx

	for rows.Next() {
		f, err = scanRowIntoIDMTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *IdmTxRepo) GetUniqueIDMTxByParentTxId(ctx context.Context, parentTxId int64) (*IDMTx, error) {
	rows, err := repo.db.Query(sQLIDMTxGetUniqueByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *IDMTx

	for rows.Next() {
		f, err = scanRowIntoIDMTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoIDMTx(rows *sql.Rows) (*IDMTx, error) {
	tx := new(IDMTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.Ip,
		&tx.UserAgent,
		&tx.IDMReference,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreatedAt,
		&tx.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func NewIDMComplianceRepo(db *sql.DB) *IDMComplianceRepo {
	return &IDMComplianceRepo{db: db}
}

type IDMComplianceRepo struct {
	db *sql.DB
}

func (repo *IDMComplianceRepo) InsertIDMComplance(ctx context.Context, c IDMCompliance) (int64, error) {
	result, err := repo.db.Exec(
		sQLIDMComplianceInsert,
		c.IDMTxId,
		c.IDMHttpStatusCode,
		c.IDMResult,
		c.RequestJson,
		c.ResponseJson,
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

func (repo *IDMComplianceRepo) GetUniqueIDMComplianceById(ctx context.Context, id int64) (*IDMCompliance, error) {
	rows, err := repo.db.Query(sQLIDMComplianceGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *IDMCompliance

	for rows.Next() {
		f, err = scanRowIntoIDMCompliance(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *IDMComplianceRepo) GetUniqueIDMComplianceByIDMTxId(ctx context.Context, idmTxId int64) (*IDMCompliance, error) {
	rows, err := repo.db.Query(sQLIDMComplianceGetUniqueByIDMTxId, idmTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *IDMCompliance

	for rows.Next() {
		f, err = scanRowIntoIDMCompliance(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoIDMCompliance(rows *sql.Rows) (*IDMCompliance, error) {
	c := new(IDMCompliance)
	err := rows.Scan(
		&c.ID,
		&c.IDMTxId,
		&c.IDMHttpStatusCode,
		&c.IDMResult,
		&c.RequestJson,
		&c.ResponseJson,
		&c.CreatedAt,
		&c.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}
