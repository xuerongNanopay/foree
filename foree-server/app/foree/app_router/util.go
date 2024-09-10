package foree_router

import (
	"context"
	"fmt"
	"net/http"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
	http_util "xue.io/go-pay/util/http"
	reflect_util "xue.io/go-pay/util/reflect"
)

func commonEndFunc[P any, Q any](req P, resp Q, hErr transport.HError) {
	if reflect_util.IsNil(hErr) {
		return
	}

	if v, is := hErr.(*transport.InteralServerError); is {
		//TODO: alert
		foree_logger.Logger.Error("System Error", "cause", v.OriginalError.Error())
	}
}

func validatePayloadBeforeProcess[P transport.ForeeRequest](r *http.Request, req P) transport.HError {
	err := req.Validate()
	if err != nil {
		foree_logger.Logger.Warn(
			"Invalid Input Error",
			"ip", http_util.LoadRealIp(r),
			"sessionId", r.Header.Get("SESSION_ID"),
			"requestType", fmt.Sprintf("%T", req),
			"cause", req.Validate().Error())
	}
	return err
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
