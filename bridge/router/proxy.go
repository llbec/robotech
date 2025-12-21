package router

import (
	"net/http"
	"net/http/httputil"
)

func NewProxyHandler(target *Service) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target.URL)
	return func(w http.ResponseWriter, r *http.Request) {
		// 可在此做统一日志 / 鉴权
		proxy.ServeHTTP(w, r)
	}
}
