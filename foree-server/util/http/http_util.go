package http_util

import (
	"net/http"
	"strings"
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
