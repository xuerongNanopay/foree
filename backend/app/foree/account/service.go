package account

import (
	"context"

	"xue.io/go-pay/app/foree/auth"
	"xue.io/go-pay/app/foree/transport"
)

type AccountService struct {
	authService *auth.AuthService
	contactRepo *ContactAccountRepo
	interacRepo *InteracAccountRepo
}

// The method is only used by CreateUser func
// So the permission check is already in there.
// We don't need permission check here.
func (a *AccountService) CreateDefaultInteracAccount(ctx context.Context, req DefaultInteracReq) transport.ForeeError {
	return nil
}
