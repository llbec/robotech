package handler

import (
	"bridge/router"
	"io"
	"net/http"
)

func MetricsHandler(reg *router.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services := reg.ListServices()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		for _, svc := range services {
			resp, err := http.Get(svc.URL.String() + svc.MetricsURL)
			if err != nil {
				continue
			}
			io.Copy(w, resp.Body)
			resp.Body.Close()
		}
	}
}
