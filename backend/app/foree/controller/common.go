package foree_controller

import (
	"context"
	"fmt"
	"net/http"

	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
)

func commonEndFunc[P any, Q any](req P, resp Q, hErr transport.HError) {
	if v, is := hErr.(*transport.InteralServerError); is {
		// use logger.
		fmt.Print(v.OriginalError.Error())
	} else {
		fmt.Println(hErr.Error())
	}
}

func commonBeforeProcess[P any](r *http.Request, req P) transport.HError {
	return nil
}

func emptyAfterProcess[Q any](w http.ResponseWriter, resp Q, hErr transport.HError) http.ResponseWriter {
	return w
}

func simpleGetWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestGetWrapper(
		handler,
		commonBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}

func simplePostWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestPostWrapper(
		handler,
		commonBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}
