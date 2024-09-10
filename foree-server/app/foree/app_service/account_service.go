package foree_service

import (
	"context"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/server/transport"
)

func NewAccountService(
	authService *AuthService,
	contactAccountRepo *account.ContactAccountRepo,
	interacAccountRepo *account.InteracAccountRepo,
) *AccountService {
	return &AccountService{
		authService:        authService,
		contactAccountRepo: contactAccountRepo,
		interacAccountRepo: interacAccountRepo,
	}
}

type AccountService struct {
	authService        *AuthService
	contactAccountRepo *account.ContactAccountRepo
	interacAccountRepo *account.InteracAccountRepo
}

func (a *AccountService) CreateContact(ctx context.Context, req CreateContactReq) (*ContactAccountDetailDTO, transport.HError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionContactWrite)
	if err != nil {
		return nil, err
	}
	now := time.Now()
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
		LatestActivityAt:      &now,
	}

	newAcc.HashMyself()

	accId, derr := a.contactAccountRepo.InsertContactAccount(ctx, newAcc)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	nAcc, nErr := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, accId)
	if nErr != nil {
		return nil, transport.WrapInteralServerError(nErr)
	}

	if nAcc == nil {
		return nil, transport.WrapInteralServerError(fmt.Errorf("can not retrieve created contact `%v`", accId))
	}

	return NewContactAccountDetailDTO(nAcc), nil
}

func (a *AccountService) DeleteContact(ctx context.Context, req DeleteContactReq) transport.HError {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionContactWrite)
	if err != nil {
		return err
	}
	acc, derr := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return transport.NewFormError("Invaild contact deletion", "contactId", "Invalid contactId")
	}

	newAcc := *acc
	newAcc.Status = account.AccountStatusDelete
	derr = a.contactAccountRepo.UpdateActiveContactAccountByIdAndOwner(ctx, newAcc)
	if derr != nil {
		return transport.WrapInteralServerError(derr)
	}
	return nil
}

func (a *AccountService) GetActiveContact(ctx context.Context, req GetContactReq) (*ContactAccountDetailDTO, transport.HError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionContactRead)
	if err != nil {
		return nil, err
	}

	acc, derr := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}

	if acc == nil {
		return nil, transport.NewFormError("Invaild contact det", "contactId", "Invalid contactId")
	}

	return NewContactAccountDetailDTO(acc), nil
}

func (a *AccountService) GetAllActiveContacts(ctx context.Context, req transport.SessionReq) ([]*ContactAccountSummaryDTO, transport.HError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionContactRead)
	if err != nil {
		return nil, err
	}

	accs, derr := a.contactAccountRepo.GetAllActiveContactAccountByOwnerId(ctx, session.User.ID)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewContactAccountSummaryDTO(v))
	}

	return ret, nil
}

func (a *AccountService) QueryActiveContacts(ctx context.Context, req QueryContactReq) ([]*ContactAccountSummaryDTO, transport.HError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionContactRead)
	if err != nil {
		return nil, err
	}
	accs, derr := a.contactAccountRepo.QueryActiveContactAccountByOwnerIdWithPagination(ctx, session.User.ID, req.Limit, req.Offset)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewContactAccountSummaryDTO(v))
	}

	return ret, nil
}

func (a *AccountService) GetAllActiveInteracs(ctx context.Context, req transport.SessionReq) ([]*InteracAccountSummaryDTO, transport.HError) {
	session, err := a.authService.Authorize(ctx, req.SessionId, PermissionInteracRead)
	if err != nil {
		return nil, err
	}

	accs, derr := a.interacAccountRepo.GetAllActiveInteracAccountByOwnerId(ctx, session.User.ID)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	ret := make([]*InteracAccountSummaryDTO, len(accs))
	for _, v := range accs {
		ret = append(ret, NewInteracAccountSummaryDTO(v))
	}

	return ret, nil
}
