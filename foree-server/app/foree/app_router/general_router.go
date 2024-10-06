package foree_router

import (
	"context"

	"github.com/gorilla/mux"
	foree_service "xue.io/go-pay/app/foree/app_service"
	"xue.io/go-pay/server/transport"
)

type GeneralRouter struct {
}

func NewGeneralRouter() *GeneralRouter {
	return &GeneralRouter{}
}

func (c *GeneralRouter) RegisterRouter(router *mux.Router) {
	// Customer Supports
	router.HandleFunc("/customer_support", simplePostWrapper(cusomterSupport)).Methods("GET")
}

func cusomterSupport(ctx context.Context, req transport.SessionReq) (*foree_service.CustomerSupport, transport.HError) {
	return &foree_service.CustomerSupport{
		SupportEmail:       "support@paypay.net",
		SupportPhoneNumber: "+1(306)555-5555",
		Instagram:          "@paypay",
		Twitter:            "@paypay",
		Facebook:           "@paypay",
	}, nil
}
