package transport

type ForeeRequest interface {
	Validate() *BadRequestError
}

type ForeeSessionRequest interface {
	Validate() *BadRequestError
	GetSessionId() string
}

type SessionReq struct {
	SessionId string
}

func (q SessionReq) GetSessionId() string {
	return q.SessionId
}

func (q SessionReq) Validate() *BadRequestError {
	return nil
}
