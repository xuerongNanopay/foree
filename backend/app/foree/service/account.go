package service

import (
	"context"
	"fmt"

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

func (a *AccountService) CreateContact(ctx context.Context, req CreateContactReq) (*ContactAccountDetailDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, ACCOUNT_CREATE)
	if err != nil {
		return nil, err
	}

	newAcc := account.ContactAccount{
		Status:                account.AccountStatusActive,
		Type:                  account.ContactAccountType(req.TransferMethod),
		FirstName:             req.FirstName,
		MiddleName:            req.MiddleName,
		LastName:              req.LastName,
		Address1:              req.Address1,
		Address2:              req.Address2,
		City:                  req.City,
		Province:              req.Province,
		Country:               req.Country,
		PostalCode:            req.PostalCode,
		PhoneNumber:           req.PhoneNumber,
		InstitutionName:       req.BankName,
		AccountNumber:         req.AccountNoOrIBAN,
		RelationshipToContact: req.RelationshipToContact,
		OwnerId:               session.User.ID,
	}

	newAcc.HashMyself()

	accId, derr := a.contactRepo.InsertContactAccount(newAcc)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	nAcc, nErr := a.contactRepo.GetUniqueContactAccountById(session.User.ID, accId)
	if nErr != nil {
		return nil, transport.WrapInteralServerError(nErr)
	}

	if nAcc == nil {
		return nil, transport.WrapInteralServerError(fmt.Errorf("can not retrieve created contact `%v`", accId))
	}

	return NewContactAccountDetailDTO(nAcc), nil
}

func (a *AccountService) DeleteContact(ctx context.Context, req DeleteContactReq) transport.ForeeError {
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

func (a *AccountService) GetContact(ctx context.Context, req GetContactReq) (*ContactAccountDetailDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, ACCOUNT_GET)
	if err != nil {
		return nil, err
	}

	acc, derr := a.contactRepo.GetUniqueContactAccountById(session.User.ID, req.ContactId)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return nil, transport.NewFormError("Invaild contact det", "contactId", "Invalid contactId")
	}

	return NewContactAccountDetailDTO(acc), nil
}

func (a *AccountService) queryContact(ctx context.Context, req QueryContactReq) ([]*ContactAccountSummaryDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, ACCOUNT_QUERY)
	if err != nil {
		return nil, err
	}
}
