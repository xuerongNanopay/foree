package promotion

import (
	"sync/atomic"
	"time"

	"xue.io/go-pay/app/foree/types"
)

type Promotion struct {
	Name      string           `json:"id"`
	Quantity  int32            `json:"limit"`
	Amt       types.AmountData `json:"Amt"`
	StartTime time.Time        `json:"startTime"`
	EndTime   time.Time        `json:"endTime"`
	IsEnable  time.Time        `json:"isEnable"`
	CreateAt  time.Time        `json:"createAt"`
	UpdateAt  time.Time        `json:"updateAt"`
}

func (p *Promotion) TryApply(name string) bool {
	if name != p.Name {
		return false
	}

	if atomic.LoadInt32(&p.Quantity) == 0 {
		return false
	}

	atomic.AddInt32(&p.Quantity, -1)

	now := time.Now().Unix()

	if now > p.StartTime.Unix() || (now > p.EndTime.Unix() && !p.EndTime.IsZero()) {
		return false
	}

	return true
}
