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

const sessionIdKey = "SESSION_ID"

func RestGetWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeProcess func(context.Context, *http.Request, P) (context.Context, transport.HError),
	afterProcess func(context.Context, http.ResponseWriter, Q, transport.HError) (context.Context, http.ResponseWriter),
	endFunc func(context.Context, P, Q, transport.HError), asyncEnd bool,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		ctx := context.Background()

		resp, herr := func() (Q, transport.HError) {
			params := mux.Vars(r)
			for k, v := range params {
				reflect_util.TrySetStuctValueFromString(&req, k, v)
			}

			for _, k := range reflect_util.GetAllFieldNamesOfStruct(&req) {
				sField, sTag := reflect_util.GetTagOfStruct(&req, k)
				rawTag, ok := sTag.Lookup("json")
				paramName := sField.Name
				if ok {
					paramName = rawTag
				}
				query := r.URL.Query()
				if query.Has(paramName) {
					reflect_util.TrySetStuctValueFromString(&req, sField.Name, query.Get(paramName))
				}
			}
			sessionId := r.Header.Get(sessionIdKey)
			reflect_util.TrySetStuctValueFromString(&req, "SessionId", sessionId)

			var resp Q
			var herr transport.HError
			ctx, herr = beforeProcess(ctx, r, req)
			if reflect_util.IsNil(herr) {
				ctx = context.WithValue(ctx, constant.CKHttpRequest, r)
				resp, herr = handler(ctx, req)
			}
			ctx, w = afterProcess(ctx, w, resp, herr)

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
			go endFunc(ctx, req, resp, herr)
		} else {
			endFunc(ctx, req, resp, herr)
		}
	}
}

func RestPostWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeProcess func(context.Context, *http.Request, P) (context.Context, transport.HError),
	afterProcess func(context.Context, http.ResponseWriter, Q, transport.HError) (context.Context, http.ResponseWriter),
	endFunc func(context.Context, P, Q, transport.HError),
	asyncEnd bool,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		ctx := context.Background()
		resp, herr := func() (Q, transport.HError) {

			sessionId := r.Header.Get(sessionIdKey)
			reflect_util.TrySetStuctValueFromString(&req, "SessionId", sessionId)

			var resp Q
			var err error
			var herr transport.HError
			err = json_util.DeserializeJsonFromHttpRequest(r, &req)
			if err != nil {
				herr = transport.WrapInteralServerError(err)
				goto SKIP_PROCESS
			}

			ctx, herr = beforeProcess(ctx, r, req)
			if reflect_util.IsNil(herr) {
				ctx := context.Background()
				ctx = context.WithValue(ctx, constant.CKHttpRequest, r)
				resp, herr = handler(ctx, req)
			}
			ctx, w = afterProcess(ctx, w, resp, herr)

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
			go endFunc(ctx, req, resp, herr)
		} else {
			endFunc(ctx, req, resp, herr)
		}
	}
}
