package service

import (
	"context"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/auth"
	"xue.io/go-pay/app/foree/transport"
)

type AccountService struct {
	authService *auth.AuthService
	contactRepo *account.ContactAccountRepo
	interacRepo *account.InteracAccountRepo
}

// The method is only used by CreateUser func
// So the permission check is already in there.
// We don't need permission check here.
func (a *AccountService) CreateDefaultInteracAccount(ctx context.Context, req account.DefaultInteracReq) transport.ForeeError {
	acc := account.InteracAccount{
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		OwnerId:     req.OwnerId,
		Status:      account.AccountStatusActive,
	}
	_, err := a.interacRepo.InsertInteracAccount(acc)
	if err != nil {
		return transport.WrapInteralServerError(err)
	}

	return nil
}
