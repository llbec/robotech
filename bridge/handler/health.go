package handler

import (
	"bridge/router"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type HealthStatus struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func HealthHandler(reg *router.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services := reg.ListServices()
		result := HealthStatus{
			Status:   "healthy",
			Services: map[string]string{},
		}

		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, svc := range services {
			wg.Add(1)
			go func(s *router.Service) {
				defer wg.Done()
				client := http.Client{Timeout: 2 * time.Second}
				resp, err := client.Get(s.URL.String() + s.HealthURL)
				status := "unhealthy"
				if err == nil && resp.StatusCode == 200 {
					status = "healthy"
				}
				mu.Lock()
				result.Services[s.Name] = status
				mu.Unlock()
			}(svc)
		}

		wg.Wait()
		// 如果任意服务 unhealthy，则总状态为 degraded
		for _, s := range result.Services {
			if s != "healthy" {
				result.Status = "degraded"
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
