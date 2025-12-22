package health

import (
	"context"
	"net/http"
	"time"

	"bridge/internal/metrics"
	"bridge/internal/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Start(reg *registry.Registry, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			reg.Iter(func(s *registry.Service) {
				switch s.Protocol {
				case "http":
					checkHTTP(s)
				case "grpc":
					checkGRPC(s)
				}
				metrics.HealthStatus.WithLabelValues(s.Name, s.Addr).Set(boolToFloat(s.Healthy.Load()))
			})
		}
	}()
}

func checkHTTP(s *registry.Service) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(s.Addr + "/health")
	if err != nil || resp.StatusCode != 200 {
		s.Healthy.Store(false)
		return
	}
	s.Healthy.Store(true)
}

func checkGRPC(s *registry.Service) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, s.Addr, grpc.WithInsecure())
	if err != nil {
		s.Healthy.Store(false)
		return
	}
	defer conn.Close()
	client := grpc_health_v1.NewHealthClient(conn)
	_, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		s.Healthy.Store(false)
		return
	}
	s.Healthy.Store(true)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
