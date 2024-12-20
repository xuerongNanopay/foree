package foree_account_service

import (
	"context"
	"fmt"

	"xue.io/go-pay/app/foree/account"
	foree_logger "xue.io/go-pay/app/foree/logger"
	foree_auth_service "xue.io/go-pay/app/foree/service/auth"
	foree_util "xue.io/go-pay/app/foree/util"
	"xue.io/go-pay/server/transport"
)

func NewAccountService(
	authService *foree_auth_service.AuthService,
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
	authService        *foree_auth_service.AuthService
	contactAccountRepo *account.ContactAccountRepo
	interacAccountRepo *account.InteracAccountRepo
}

func (a *AccountService) VerifyContact(ctx context.Context, req CreateContactReq) (*VerifyContactDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("VerifyContact_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	//TODO: user real services.
	// var resp *nbp.AccountEnquiryResponse
	// if req.AccountNoOrIBAN == "1111" {
	// 	resp = &nbp.AccountEnquiryResponse{
	// 		ResponseCommon: nbp.ResponseCommon{
	// 			ResponseCode: "407",
	// 		},
	// 	}
	// } else {
	// 	resp = &nbp.AccountEnquiryResponse{
	// 		ResponseCommon: nbp.ResponseCommon{
	// 			ResponseCode: "201",
	// 		},
	// 	}
	// }

	if req.AccountNoOrIBAN == "1111" {
		return &VerifyContactDTO{
			AccountStatus: "Closed",
		}, nil
	} else if req.AccountNoOrIBAN == "4444" {
		return &VerifyContactDTO{
			AccountStatus: "BUSINESS",
		}, nil
	} else if req.AccountNoOrIBAN == "5555" {
		return &VerifyContactDTO{}, nil
	} else {
		return &VerifyContactDTO{
			AccountStatus: "Active",
		}, nil
	}
}

func (a *AccountService) CreateContact(ctx context.Context, req CreateContactReq) (*ContactAccountDetailDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CreateContact_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
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

	accId, err := a.contactAccountRepo.InsertContactAccount(ctx, newAcc)
	if err != nil {
		foree_logger.Logger.Error("CreateContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	nAcc, err := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, accId)
	if err != nil {
		foree_logger.Logger.Error("CreateContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if nAcc == nil {
		foree_logger.Logger.Error("CreateContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", fmt.Errorf("can not retrieve created contact with id `%v`", accId))
		return nil, transport.WrapInteralServerError(fmt.Errorf("can not retrieve created contact with id `%v`", accId))
	}

	foree_logger.Logger.Info("CreateContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "contactId", accId)
	return NewContactAccountDetailDTO(nAcc), nil
}

func (a *AccountService) DeleteContact(ctx context.Context, req DeleteContactReq) (*ContactAccountDetailDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("DeleteContact_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	acc, err := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if err != nil {
		foree_logger.Logger.Error("DeleteContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if acc == nil {
		foree_logger.Logger.Error("DeleteContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", fmt.Sprintf("can't find active contact with id `%v`", req.ContactId))
		return nil, transport.NewFormError("Invaild contact deletion", "contactId", "Invalid contactId")
	}

	newAcc := *acc
	newAcc.Status = account.AccountStatusDelete
	err = a.contactAccountRepo.UpdateActiveContactAccountByIdAndOwner(ctx, newAcc)
	if err != nil {
		foree_logger.Logger.Error("DeleteContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("DeleteContact_SUCCESS", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "contactId", req.ContactId)
	return nil, nil
}

func (a *AccountService) GetActiveContact(ctx context.Context, req GetContactReq) (*ContactAccountDetailDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetActiveContact_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	acc, err := a.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, session.User.ID, req.ContactId)
	if err != nil {
		foree_logger.Logger.Error("GetActiveContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if acc == nil {
		foree_logger.Logger.Error("GetActiveContact_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", fmt.Sprintf("can't find active contact with id `%v`", req.ContactId))
		return nil, transport.NewFormError("Invaild contact det", "contactId", "Invalid contactId")
	}

	foree_logger.Logger.Debug("GetActiveContact_SUCCESS", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId)
	return NewContactAccountDetailDTO(acc), nil
}

func (a *AccountService) GetAllActiveContacts(ctx context.Context, req transport.SessionReq) ([]*ContactAccountSummaryDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetAllActiveContacts_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	accs, err := a.contactAccountRepo.GetAllActiveContactAccountByOwnerId(ctx, session.User.ID)
	if err != nil {
		foree_logger.Logger.Error("GetAllActiveContacts_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for i, v := range accs {
		ret[i] = NewContactAccountSummaryDTO(v)
	}
	foree_logger.Logger.Debug("GetAllActiveContacts_SUCCESS", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId)
	return ret, nil
}

func (a *AccountService) QueryActiveContacts(ctx context.Context, req QueryContactReq) ([]*ContactAccountSummaryDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("QueryActiveContacts_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	accs, err := a.contactAccountRepo.QueryActiveContactAccountByOwnerIdWithPagination(ctx, session.User.ID, req.Limit, req.Offset)
	if err != nil {
		foree_logger.Logger.Error("QueryActiveContacts_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	ret := make([]*ContactAccountSummaryDTO, len(accs))
	for i, v := range accs {
		ret[i] = NewContactAccountSummaryDTO(v)
	}

	foree_logger.Logger.Debug("QueryActiveContacts_SUCCESS", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId)
	return ret, nil
}

func (a *AccountService) GetAllActiveInteracs(ctx context.Context, req transport.SessionReq) ([]*InteracAccountSummaryDTO, transport.HError) {
	session, sErr := a.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetAllActiveInteracs_FAIL", "sessionId", req.SessionId, "cause", "session no found")
		return nil, sErr
	}

	accs, err := a.interacAccountRepo.GetAllActiveInteracAccountByOwnerId(ctx, session.User.ID)
	if err != nil {
		foree_logger.Logger.Error("GetAllActiveInteracs_FAIL", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	ret := make([]*InteracAccountSummaryDTO, len(accs))
	for i, v := range accs {
		ret[i] = NewInteracAccountSummaryDTO(v)
	}

	foree_logger.Logger.Debug("GetAllActiveInteracs_SUCCESS", "ip", foree_util.LoadRealIp(ctx), "userId", session.UserId)
	return ret, nil
}
