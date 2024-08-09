package transaction

import (
	"fmt"
	"math"

	"xue.io/go-pay/app/foree/types"
)

type Rate struct {
	SrcAmt  types.AmountData `json:"srcAmt"`
	DestAmt types.AmountData `json:"destAmt"`
}

func (r *Rate) ToSummary() string {
	return fmt.Sprintf("$%.2f %s : %.2f %s", r.SrcAmt.Amount, r.SrcAmt.Curreny, r.DestAmt.Amount, r.DestAmt.Curreny)
}

func (r *Rate) GetForwardRate() float64 {
	return math.Round((float64(r.DestAmt.Amount)/float64(r.SrcAmt.Amount))*100) / 100
}

func (r *Rate) GetBackwardRate() float64 {
	return math.Round((float64(r.SrcAmt.Amount)/float64(r.DestAmt.Amount))*100) / 100
}
