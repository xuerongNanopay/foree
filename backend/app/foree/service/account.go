package service

import (
	"context"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transport"
)

type AccountService struct {
	authService *AuthService
	contactRepo *account.ContactAccountRepo
	interacRepo *account.InteracAccountRepo
}

func (a *AccountService) CreateContact(ctx context.Context, req CreateContactReq) (*ContactAccountDetailDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, Contact_CREATE)
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
		LatestActivityAt:      time.Now(),
	}

	newAcc.HashMyself()

	accId, derr := a.contactRepo.InsertContactAccount(ctx, newAcc)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	nAcc, nErr := a.contactRepo.GetUniqueNonDeleteContactAccountByOwnerAndId(ctx, session.User.ID, accId)
	if nErr != nil {
		return nil, transport.WrapInteralServerError(nErr)
	}

	if nAcc == nil {
		return nil, transport.WrapInteralServerError(fmt.Errorf("can not retrieve created contact `%v`", accId))
	}

	return NewContactAccountDetailDTO(nAcc), nil
}

func (a *AccountService) DeleteContact(ctx context.Context, req DeleteContactReq) transport.ForeeError {
	session, err := a.authService.Authorize(ctx, req.SessionId, Contact_DELETE)
	if err != nil {
		return err
	}
	acc, derr := a.contactRepo.GetUniqueNonDeleteContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return transport.NewFormError("Invaild contact deletion", "contactId", "Invalid contactId")
	}

	newAcc := *acc
	newAcc.Status = account.AccountStatusDelete
	derr = a.contactRepo.UpdateNonDeleteContactAccountByIdAndOwner(ctx, newAcc)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}
	return nil
}

func (a *AccountService) GetContact(ctx context.Context, req GetContactReq) (*ContactAccountDetailDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, Contact_GET)
	if err != nil {
		return nil, err
	}

	acc, derr := a.contactRepo.GetUniqueNonDeleteContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return nil, transport.NewFormError("Invaild contact det", "contactId", "Invalid contactId")
	}

	return NewContactAccountDetailDTO(acc), nil
}

func (a *AccountService) GetAllContacts(ctx context.Context, req transport.SessionReq) ([]*ContactAccountSummaryDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, Contact_QUERY)
	if err != nil {
		return nil, err
	}

	accs, derr := a.contactRepo.GetAllNonDeleteContactAccountByOwnerId(ctx, session.User.ID)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewContactAccountSummaryDTO(v))
	}

	return ret, nil
}

func (a *AccountService) QueryContacts(ctx context.Context, req QueryContactReq) ([]*ContactAccountSummaryDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, Contact_QUERY)
	if err != nil {
		return nil, err
	}
	accs, derr := a.contactRepo.QueryNonDeleteContactAccountByOwnerId(ctx, session.User.ID, req.Limit, req.Offset)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewContactAccountSummaryDTO(v))
	}

	return ret, nil
}

func (a *AccountService) GetAllInteracs(ctx context.Context, req transport.SessionReq) ([]*InteracAccountSummaryDTO, transport.ForeeError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermInteracQuery)
	if err != nil {
		return nil, err
	}

	accs, derr := a.interacRepo.GetAllNonDeleteInteracAccountByOwnerId(ctx, session.User.ID)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*InteracAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewInteracAccountSummaryDTO(v))
	}

	return ret, nil
}
