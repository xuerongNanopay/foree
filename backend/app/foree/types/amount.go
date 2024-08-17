package types

import (
	"fmt"
	"strings"
)

type Amount float64

func (a Amount) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", a)
	return []byte(s), nil
}

type AmountData struct {
	Amount   Amount `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type RateDate struct {
	Src  AmountData
	Dest AmountData
}

func (r RateDate) String() string {
	return fmt.Sprintf("%.2f%v:%.2f%v", r.Src.Amount, strings.ToUpper(r.Src.Currency), r.Dest.Amount, strings.ToUpper(r.Dest.Currency))
}
