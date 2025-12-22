package httpx

import (
	"fmt"
	"net/http"
	"strconv"

	"bridge/internal/limiter"
	"bridge/internal/registry"
)

func NewHandler(reg *registry.Registry) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("use /metrics for Prometheus"))
	})

	mux.HandleFunc("/admin/register", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		addr := r.URL.Query().Get("addr")
		proto := r.URL.Query().Get("proto")
		weight, _ := strconv.Atoi(r.URL.Query().Get("weight"))
		gray, _ := strconv.Atoi(r.URL.Query().Get("gray"))
		qps, _ := strconv.Atoi(r.URL.Query().Get("qps"))
		if qps == 0 {
			qps = 50
		}
		if weight == 0 {
			weight = 100
		}
		s := &registry.Service{
			Name:        name,
			Addr:        addr,
			Protocol:    proto,
			Weight:      weight,
			GrayPercent: gray,
		}
		reg.Register(s, qps)
		fmt.Fprintf(w, "registered %s %s\n", name, addr)
	})

	mux.HandleFunc("/admin/update", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		addr := r.URL.Query().Get("addr")
		weight, _ := strconv.Atoi(r.URL.Query().Get("weight"))
		gray, _ := strconv.Atoi(r.URL.Query().Get("gray"))
		qps, _ := strconv.Atoi(r.URL.Query().Get("qps"))
		found := false
		reg.Iter(func(s *registry.Service) {
			if s.Name == name && s.Addr == addr {
				if weight > 0 {
					s.Weight = weight
				}
				if gray >= 0 {
					s.GrayPercent = gray
				}
				if qps > 0 {
					s.Limiter = limiter.NewLimiter(qps, qps)
				}
				found = true
			}
		})
		if found {
			fmt.Fprintf(w, "updated %s %s\n", name, addr)
		} else {
			http.Error(w, "service not found", 404)
		}
	})

	return mux
}
