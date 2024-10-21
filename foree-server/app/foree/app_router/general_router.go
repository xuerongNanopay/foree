package foree_router

import (
	"context"

	"github.com/gorilla/mux"
	"xue.io/go-pay/server/transport"
)

type CustomerSupport struct {
	SupportEmail       string `json:"supportEmail"`
	SupportPhoneNumber string `json:"supportPhoneNumber"`
	Instagram          string `json:"instagram"`
	Twitter            string `json:"twitter"`
	Facebook           string `json:"facebook"`
}

type GeneralRouter struct {
}

func NewGeneralRouter() *GeneralRouter {
	return &GeneralRouter{}
}

func (c *GeneralRouter) RegisterRouter(router *mux.Router) {
	// Customer Supports
	router.HandleFunc("/customer_support", simpleGetWrapper(cusomterSupport)).Methods("GET")
}

func cusomterSupport(ctx context.Context, req transport.SessionReq) (*CustomerSupport, transport.HError) {
	return &CustomerSupport{
		SupportEmail:       "support@paypay.net",
		SupportPhoneNumber: "+1(306)555-5555",
		Instagram:          "@paypay",
		Twitter:            "@paypay",
		Facebook:           "@paypay",
	}, nil
}
