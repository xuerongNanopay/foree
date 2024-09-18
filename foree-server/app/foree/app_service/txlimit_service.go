package foree_service

import (
	"context"
	"fmt"

	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
)

func NewTxLimitService(dailyTxLimiteRepo *transaction.DailyTxLimitRepo) *TxLimitService {
	return &TxLimitService{
		dailyTxLimiteRepo: dailyTxLimiteRepo,
	}
}

type TxLimitService struct {
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo
}

func (t *TxLimitService) addDailyTxLimit(ctx context.Context, session auth.Session, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := t.getDailyTxLimit(ctx, session)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount += amt.Amount

	if err := t.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}

	return dailyLimit, nil
}

// I don't case race condition here, cause create transaction will rescure it.
func (t *TxLimitService) getDailyTxLimit(ctx context.Context, session auth.Session) (*transaction.DailyTxLimit, error) {
	reference := transaction.GenerateDailyTxLimitReference(session.UserId)
	dailyLimit, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
	if err != nil {
		return nil, err
	}

	// If not create one.
	if dailyLimit == nil {
		txLimit, ok := foree_constant.TxLimits[foree_constant.TransactionLimitGroup(session.UserGroup.TransactionLimitGroup)]
		if !ok {
			return nil, fmt.Errorf("transaction limit no found for group `%v`", session.UserGroup.TransactionLimitGroup)
		}
		dailyLimit = &transaction.DailyTxLimit{
			Reference: reference,
			UsedAmt: types.AmountData{
				Amount:   0.0,
				Currency: txLimit.MaxAmt.Currency,
			},
			MaxAmt: types.AmountData{
				Amount:   txLimit.MaxAmt.Amount,
				Currency: txLimit.MaxAmt.Currency,
			},
			OwnerId: session.UserId,
		}
		_, err := t.dailyTxLimiteRepo.InsertDailyTxLimit(ctx, *dailyLimit)
		if err != nil {
			return nil, err
		}
		dl, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
		if err != nil {
			return nil, err
		}
		dailyLimit = dl
	}
	return dailyLimit, nil
}
