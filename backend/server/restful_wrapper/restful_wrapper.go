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

func RestPostWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	afterHandler func(P, Q, transport.HError), isAsyncAfter bool,
	customeWriter func(http.ResponseWriter) http.ResponseWriter,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		resp, herr := func() (Q, transport.HError) {
			if err := json_util.DeserializeJsonFromHttpRequest(r, &req); err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			}

			sessionId := r.Header.Get("SESSION_ID")
			reflect_util.SetStringValueIfFieldExist(&req, "SessionId", sessionId)

			return handler(context.Background(), req)
		}()

		if isAsyncAfter {
			go afterHandler(req, resp, herr)
		} else {
			afterHandler(req, resp, herr)
		}

		w.Header().Add("Content-Type", "application/json")
		w = customeWriter(w)

		var err error
		if herr != nil {
			err = json_util.SerializeToResponseWriter(w, herr.GetStatusCode(), herr)
		} else {
			err = json_util.SerializeToResponseWriter(w, http.StatusOK, transport.NewHttpResponse(http.StatusOK, "Success", resp))
		}

		if err != nil {
			//TODO: log error
			return
		}
	}
}
