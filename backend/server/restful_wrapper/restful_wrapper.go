package restful_wrapper

import (
	"context"
	"net/http"

	"xue.io/go-pay/server/transport"
)

func RestGetWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func RestPostWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func RestPutWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func RestDeleteWrapper[P any, Q any](f func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
