package transport

type ForeeRequest interface {
	TrimSpace()
	Validate() *BadRequestError
}

type SessionReq struct {
	SessionId string
}

type ISession interface {
	SetSession(string)
	GetSession() string
}

func (s *SessionReq) SetSession(sessionId string) {
	s.SessionId = sessionId
}

func (s *SessionReq) GetSession() string {
	return s.SessionId
}
