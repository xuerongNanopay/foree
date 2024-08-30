package restful_wrapper

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"xue.io/go-pay/server/transport"
	json_util "xue.io/go-pay/util/json"
	reflect_util "xue.io/go-pay/util/reflect"
)

func RestGetWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeResponse func(http.ResponseWriter) http.ResponseWriter,
	afterRun func(P, Q, transport.HError), isAsyncAfter bool,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req P
		resp, herr := func() (Q, transport.HError) {
			params := mux.Vars(r)
			for k, v := range params {
				reflect_util.SetIntOrStringValueIfFieldExistFromString(&req, k, v)
			}

			for _, k := range reflect_util.GetAllFieldNamesOfStruct(&req) {
				query := r.URL.Query()
				if query.Has(k) {
					reflect_util.SetIntOrStringValueIfFieldExistFromString(&req, k, query.Get(k))
				}
			}

			sessionId := r.Header.Get("SESSION_ID")
			reflect_util.SetStringValueIfFieldExist(&req, "SessionId", sessionId)

			resp, herr := handler(context.Background(), req)

			w.Header().Add("Content-Type", "application/json")
			w = beforeResponse(w)

			var err error
			if herr != nil {
				err = json_util.SerializeToResponseWriter(w, herr.GetStatusCode(), herr)
			} else {
				err = json_util.SerializeToResponseWriter(w, http.StatusOK, transport.NewHttpResponse(http.StatusOK, "Success", resp))
			}

			if herr != nil {
				var nilResp Q
				return nilResp, herr
			} else if err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			} else {
				return resp, nil
			}

		}()

		if isAsyncAfter {
			go afterRun(req, resp, herr)
		} else {
			afterRun(req, resp, herr)
		}
	}
}

func RestPostWrapper[P any, Q any](
	handler func(context.Context, P) (Q, transport.HError),
	beforeResponse func(http.ResponseWriter) http.ResponseWriter,
	afterRun func(P, Q, transport.HError),
	isAsyncAfter bool,
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

			resp, herr := handler(context.Background(), req)

			w.Header().Add("Content-Type", "application/json")
			w = beforeResponse(w)

			var err error
			if herr != nil {
				err = json_util.SerializeToResponseWriter(w, herr.GetStatusCode(), herr)
			} else {
				err = json_util.SerializeToResponseWriter(w, http.StatusOK, transport.NewHttpResponse(http.StatusOK, "Success", resp))
			}

			if herr != nil {
				var nilResp Q
				return nilResp, herr
			} else if err != nil {
				var nilResp Q
				return nilResp, transport.WrapInteralServerError(err)
			} else {
				return resp, nil
			}

		}()

		if isAsyncAfter {
			go afterRun(req, resp, herr)
		} else {
			afterRun(req, resp, herr)
		}
	}
}
