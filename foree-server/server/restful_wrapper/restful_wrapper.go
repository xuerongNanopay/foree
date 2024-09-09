package restful_wrapper

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
	json_util "xue.io/go-pay/util/json"
	reflect_util "xue.io/go-pay/util/reflect"
)

func RestGetWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeProcess func(*http.Request, P) transport.HError,
	afterProcess func(http.ResponseWriter, Q, transport.HError) http.ResponseWriter,
	endFunc func(P, Q, transport.HError), asyncEnd bool,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		resp, herr := func() (Q, transport.HError) {
			params := mux.Vars(r)
			for k, v := range params {
				reflect_util.TrySetStuctValueFromString(&req, k, v)
			}

			for _, k := range reflect_util.GetAllFieldNamesOfStruct(&req) {
				query := r.URL.Query()
				if query.Has(k) {
					reflect_util.TrySetStuctValueFromString(&req, k, query.Get(k))
				}
			}

			sessionId := r.Header.Get("SESSION_ID")
			reflect_util.TrySetStuctValueFromString(&req, "SessionId", sessionId)

			var resp Q
			var herr transport.HError
			herr = beforeProcess(r, req)
			if reflect_util.IsNil(herr) {
				ctx := context.Background()
				ctx = context.WithValue(ctx, constant.CKHttpRequest, r)
				resp, herr = handler(ctx, req)
			}
			w = afterProcess(w, resp, herr)

			w.Header().Add("Content-Type", "application/json")

			var err error
			if !reflect_util.IsNil(herr) {
				err = json_util.SerializeToResponseWriter(w, herr.GetStatusCode(), herr)
			} else {
				err = json_util.SerializeToResponseWriter(w, http.StatusOK, transport.NewHttpResponse(http.StatusOK, "Success", resp))
			}

			if !reflect_util.IsNil(herr) {
				var nilResp Q
				return nilResp, herr
			} else if err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			} else {
				return resp, nil
			}

		}()

		if asyncEnd {
			go endFunc(req, resp, herr)
		} else {
			endFunc(req, resp, herr)
		}
	}
}

func RestPostWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeProcess func(*http.Request, P) transport.HError,
	afterProcess func(http.ResponseWriter, Q, transport.HError) http.ResponseWriter,
	endFunc func(P, Q, transport.HError),
	asyncEnd bool,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		resp, herr := func() (Q, transport.HError) {

			sessionId := r.Header.Get("SESSION_ID")
			reflect_util.TrySetStuctValueFromString(&req, "SessionId", sessionId)

			var resp Q
			var err error
			var herr transport.HError
			err = json_util.DeserializeJsonFromHttpRequest(r, &req)
			if err != nil {
				herr = transport.WrapInteralServerError(err)
				goto SKIP_PROCESS
			}

			herr = beforeProcess(r, req)
			if reflect_util.IsNil(herr) {
				ctx := context.Background()
				ctx = context.WithValue(ctx, constant.CKHttpRequest, r)
				resp, herr = handler(ctx, req)
			}
			w = afterProcess(w, resp, herr)

			w.Header().Add("Content-Type", "application/json")

		SKIP_PROCESS:
			if !reflect_util.IsNil(herr) {
				err = json_util.SerializeToResponseWriter(w, herr.GetStatusCode(), herr)
			} else {
				err = json_util.SerializeToResponseWriter(w, http.StatusOK, transport.NewHttpResponse(http.StatusOK, "Success", resp))
			}

			if !reflect_util.IsNil(herr) {
				var nilResp Q
				return nilResp, herr
			} else if err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			} else {
				return resp, nil
			}

		}()

		if asyncEnd {
			go endFunc(req, resp, herr)
		} else {
			endFunc(req, resp, herr)
		}
	}
}
