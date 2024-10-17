package nbp

import (
	"fmt"
	"net/http"
	"sync"
)

func NewMockNBPClient() NBPClient {
	return &NBPClientMock{
		config: NBPConfig{
			Mode: "mock",
		},
	}
}

type NBPClientMock struct {
	config                NBPConfig
	cancelledTransactions sync.Map
}

func (s *NBPClientMock) GetConfigs() NBPConfig {
	return s.config
}

func (s *NBPClientMock) SetConfig(key string, value string) error {
	return nil
}

func (*NBPClientMock) Hello() (*HelloResponse, error) {
	ret := &HelloResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  http.StatusOK,
			RawResponse: "Welcome, NBP E-Remittance API",
		},
	}
	return ret, nil
}

func (*NBPClientMock) BankList() (*BankListResponse, error) {
	return &BankListResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Bank list retrieved successfully",
		},
		Banklist: []BankListEntry{
			{
				BankName: "ALLIED BANK LIMITED",
				BankAbbr: "ABL",
			},
		},
	}, nil
}

func (*NBPClientMock) AccountEnquiry(r AccountEnquiryRequest) (*AccountEnquiryResponse, error) {
	return &AccountEnquiryResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Account Details retrieved successfully",
		},
		BankName:      r.BankName,
		BranchCode:    r.BranchCode,
		AccountNo:     r.AccountNo,
		AccountStatus: AccStatusActive,
	}, nil
}

func (*NBPClientMock) LoadRemittanceCash(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	return &LoadRemittanceResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Remittance credited successfully.",
		},
		GlobalId:   r.GlobalId,
		TrackingId: r.GlobalId,
	}, nil
}

func (*NBPClientMock) LoadRemittanceAccounts(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	return &LoadRemittanceResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Remittance credited successfully.",
		},
		GlobalId:   r.GlobalId,
		TrackingId: r.GlobalId,
	}, nil
}

func (*NBPClientMock) LoadRemittanceThirdParty(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	return &LoadRemittanceResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Remittance credited successfully.",
		},
		GlobalId:   r.GlobalId,
		TrackingId: r.GlobalId,
	}, nil
}

func (m *NBPClientMock) TransactionStatusByIds(r TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error) {
	var txStatuses []TransactionStatus
	ids := []string(r.Ids)

	for i := 0; i < len(ids); i++ {
		_, ok := m.cancelledTransactions.LoadAndDelete(ids[i])
		if ok {
			txStatuses = append(txStatuses, TransactionStatus{
				GlobalId: ids[i],
				Status:   TxStatusCancelled,
			})
		} else {
			txStatuses = append(txStatuses, TransactionStatus{
				GlobalId: ids[i],
				Status:   TxStatusPaid,
			})
		}
	}

	return &TransactionStatusByIdsResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Remittance status retrieved successfully",
		},
		TransactionStatuses: txStatuses,
	}, nil
}

func (*NBPClientMock) TransactionStatusByDate(r TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error) {
	return nil, fmt.Errorf("TransactionStatusByDate is not implemented for mock client")
}
func (m *NBPClientMock) CancelTransaction(r CancelTransactionRequest) (*CancelTransactionResponse, error) {
	m.cancelledTransactions.Store(r.GlobalId, true)
	return &CancelTransactionResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:      200,
			ResponseCode:    "201",
			ResponseMessage: "Remittance cancelled successfully",
		},
		GlobalId: r.GlobalId,
	}, nil
}
