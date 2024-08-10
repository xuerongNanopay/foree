package auth

import (
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/auth"
)

func verifySession(session *auth.Session) transport.ForeeError {
	if session == nil || session.EmailPasswd == nil {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgLogin,
			transport.RequireActionLogin,
		)
	}
	if session.EmailPasswd.Status == auth.EPStatusWaitingVerify {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgVerifyEmail,
			transport.RequireActionVerifyEmail,
		)
	}
	if session.User == nil || session.User.Status == auth.UserStatusInitial {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgCreateUser,
			transport.RequireActionCreateUser,
		)
	}
	return nil
}
