package idm_compliance

import (
	"database/sql"
	"time"
)

const (
	sQLIDMComplianceInsert = `
		INSERT INTO idm_compliance
		(
			foree_tx_id, idm_resp_http_code,
			idm_resp_status, idm_raw_request,
			idm_raw_response
		) VALUES(?,?,?,?,?)
	`
	sQLIDMComplianceGetUniqueById = `
		SELECT
			i.id, i.foree_tx_id, i.idm_resp_http_code,
			i.idm_resp_status, i.idm_raw_request,
			i.idm_raw_response, i.create_at, i.update_at
		FROM idm_compliance as i
		WHERE i.id = ?
	`
	sQLIDMComplianceGetUniqueByForeeTxId = `
		SELECT
			i.id, i.foree_tx_id, i.idm_resp_http_code,
			i.idm_resp_status, i.idm_raw_request,
			i.idm_raw_response, i.create_at, i.update_at
		FROM idm_compliance as i
		WHERE i.foree_tx_id = ?
	`
)

type IDMCompliance struct {
	ID              int64     `json:"id"`
	ForeeTxId       int64     `json:"foreeTxId"`
	IDMRespHttpCode int       `json:"idmRespHttpCode"`
	IDMRespStatus   string    `json:"idmRespStatus"`
	IDMRawRequest   string    `json:"idmRawRequest"`
	IDMRawResponse  string    `json:"idmRawResponse"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
}

func NewIDMComplianceRepo(db *sql.DB) *IDMComplianceRepo {
	return &IDMComplianceRepo{db: db}
}

type IDMComplianceRepo struct {
	db *sql.DB
}

func (repo *IDMComplianceRepo) InsertIDMCompliance(compliance IDMCompliance) (int64, error) {
	result, err := repo.db.Exec(
		sQLIDMComplianceInsert,
		compliance.ForeeTxId,
		compliance.IDMRespHttpCode,
		compliance.IDMRespStatus,
		compliance.IDMRawRequest,
		compliance.IDMRawResponse,
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

func (repo *IDMComplianceRepo) GetUniqueIDMComplianceById(id int64) (*IDMCompliance, error) {
	rows, err := repo.db.Query(sQLIDMComplianceGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *IDMCompliance

	for rows.Next() {
		u, err = scanRowIntoIDMCompliance(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func (repo *IDMComplianceRepo) GetUniqueIDMComplianceByForeeTxId(foreeTxId int64) (*IDMCompliance, error) {
	rows, err := repo.db.Query(sQLIDMComplianceGetUniqueByForeeTxId, foreeTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *IDMCompliance

	for rows.Next() {
		u, err = scanRowIntoIDMCompliance(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func scanRowIntoIDMCompliance(rows *sql.Rows) (*IDMCompliance, error) {
	i := new(IDMCompliance)
	err := rows.Scan(
		&i.ID,
		&i.ForeeTxId,
		&i.IDMRespHttpCode,
		&i.IDMRespStatus,
		&i.IDMRawRequest,
		&i.IDMRawResponse,
		&i.CreateAt,
		&i.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
