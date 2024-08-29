package restful_wrapper

import (
	"context"

	"xue.io/go-pay/server/transport"
)

func RestGetWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError)) {

}
