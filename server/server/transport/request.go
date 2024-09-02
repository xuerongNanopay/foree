package transport

type ForeeRequest interface {
	Validate() *BadRequestError
}

type SessionReq struct {
	SessionId string
}

func (q SessionReq) Validate() *BadRequestError {
	return nil
}
