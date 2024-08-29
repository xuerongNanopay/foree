package restful_wrapper

import (
	"context"
	"net/http"

	"xue.io/go-pay/server/transport"
	json_util "xue.io/go-pay/util/json"
	reflect_util "xue.io/go-pay/util/reflect"
)

func RestGetWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError), afterHandler func(Q, transport.HError), isAsyncAfter bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func RestPostWrapper[P any, Q any](handler func(context.Context, P) (Q, transport.HError), afterHandler func(Q, transport.HError), isAsyncAfter bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, herr := func() (Q, transport.HError) {
			var req P
			if err := json_util.ParseJsonFromHttpRequest(r, &req); err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			}

			//TODO: get session and inject into req
			sessionId := r.Header.Get("SESSION_ID")
			reflect_util.SetStringValueIfFieldExist(&req, "SessionId", sessionId)

			return handler(context.Background(), req)
		}()

		if isAsyncAfter {
			go afterHandler(resp, herr)
		} else {
			afterHandler(resp, herr)
		}
	}
}
