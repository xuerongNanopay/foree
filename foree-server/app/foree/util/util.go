package foree_util

import (
	"context"
	"net/http"

	"xue.io/go-pay/constant"
	http_util "xue.io/go-pay/util/http"
)

func LoadRealIp(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return http_util.LoadRealIp(req)
}

func LoadUserAgent(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return req.Header.Get("User-Agent")
}
