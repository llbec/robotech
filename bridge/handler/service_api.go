package handler

import (
	"bridge/router"
	"encoding/json"
	"net/http"
)

type AddServiceRequest struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	HealthURL  string `json:"health_url"`
	MetricsURL string `json:"metrics_url"`
}

func AddServiceHandler(reg *router.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddServiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := reg.AddService(req.Name, req.URL, req.HealthURL, req.MetricsURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
