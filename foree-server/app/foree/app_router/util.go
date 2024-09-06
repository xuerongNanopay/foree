package foree_router

import (
	"context"
	"net/http"

	"xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
)

func commonEndFunc[P any, Q any](req P, resp Q, hErr transport.HError) {
	if v, is := hErr.(*transport.InteralServerError); is {
		// use logger.
		logger.Logger.Error("System Error", "cause", v.OriginalError.Error())
	} else {
		logger.Logger.Warn("API Warning", "case", v.Error())
	}
}

func validatePayloadBeforeProcess[P transport.ForeeRequest](r *http.Request, req P) transport.HError {
	return req.Validate()
}

func emptyAfterProcess[Q any](w http.ResponseWriter, resp Q, hErr transport.HError) http.ResponseWriter {
	return w
}

func simpleGetWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestGetWrapper(
		handler,
		validatePayloadBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}

func simplePostWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestPostWrapper(
		handler,
		validatePayloadBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}
