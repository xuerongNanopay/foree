package service

import (
	"context"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transport"
)

type AccountService struct {
	authService *AuthService
	contactRepo *account.ContactAccountRepo
	interacRepo *account.InteracAccountRepo
}

// The method is only used by CreateUser func
// So the permission check is already in there.
// We don't need permission check here.
func (a *AccountService) CreateDefaultInteracAccount(ctx context.Context, req DefaultInteracReq) transport.ForeeError {
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

func (a *AccountService) DeleteAccount(ctx context.Context, req DeleteContactReq) transport.ForeeError {
	session, err := a.authService.Authorize(ctx, req.SessionId, ACCOUNT_DELETE)
	if err != nil {
		return err
	}
	acc, derr := a.contactRepo.GetUniqueContactAccountById(session.User.ID, req.ContactId)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return transport.NewFormError("Invaild contact deletion", "contactId", "Invalid contactId")
	}

	newAcc := *acc
	newAcc.Status = account.AccountStatusDelete
	derr = a.contactRepo.UpdateContactAccountById(newAcc)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}
	return nil
}
