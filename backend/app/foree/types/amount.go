package types

import (
	"fmt"
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
