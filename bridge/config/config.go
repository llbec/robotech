package config

type ServiceConfig struct {
	Name       string
	URL        string
	HealthURL  string
	MetricsURL string
}

type Config struct {
	Port     string
	Services []ServiceConfig
}

func DefaultConfig() *Config {
	return &Config{
		Port: "8080",
		Services: []ServiceConfig{
			{
				Name:       "txstore",
				URL:        "http://localhost:8081",
				HealthURL:  "/health",
				MetricsURL: "/metrics",
			},
			// 可以继续添加 state/report 服务
		},
	}
}
