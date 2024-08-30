package foree_controller

import (
	"context"
	"fmt"
	"net/http"

	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
)

func afterLogger[P any, Q any](req P, resp Q, hErr transport.HError) {
	if v, is := hErr.(*transport.InteralServerError); is {
		// use logger.
		fmt.Print(v.OriginalError.Error())
	} else {
		fmt.Println(hErr.Error())
	}
}

func emptyBeforeResponse[Q any](w http.ResponseWriter, resp Q) http.ResponseWriter {
	return w
}

func simpleGetWrapper[P any, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestGetWrapper(
		handler,
		emptyBeforeResponse,
		afterLogger[P, Q],
		true,
	)
}

func simplePostWrapper[P any, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestPostWrapper(
		handler,
		emptyBeforeResponse,
		afterLogger[P, Q],
		true,
	)
}
