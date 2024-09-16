package transport

type ForeeRequest interface {
	Validate() *BadRequestError
}

type SessionReq struct {
	SessionId string
}

type ForeeSessionRequest interface {
	Validate() *BadRequestError
	GetSessionId() string
}

func (q SessionReq) Validate() *BadRequestError {
	return nil
}
