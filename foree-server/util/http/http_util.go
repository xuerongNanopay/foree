package http_util

import (
	"context"
	"net/http"
	"strings"

	"xue.io/go-pay/constant"
)

func LoadRealIp(r *http.Request) string {
	xforward := r.Header.Get("X-Forwarded-For")
	var ip string
	if xforward == "" || len(strings.Split(xforward, ",")) == 0 {
		ip = r.RemoteAddr
	} else {
		ip = strings.Split(xforward, ",")[0]
	}

	if len(strings.Split(ip, ":")) > 1 {
		return strings.Split(ip, ":")[0]
	} else {
		return ip
	}
}

func LoadRealIpFromContext(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return LoadRealIp(req)
}
