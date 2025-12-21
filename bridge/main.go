package main

import (
	"bridge/config"
	"bridge/handler"
	"bridge/router"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.DefaultConfig()
	reg := router.NewRegistry()

	// 初始化服务注册表
	for _, svc := range cfg.Services {
		if err := reg.AddService(svc.Name, svc.URL, svc.HealthURL, svc.MetricsURL); err != nil {
			log.Fatalf("failed to add service: %v", err)
		}
	}

	r := chi.NewRouter()

	// 内部接口
	r.Get("/health", handler.HealthHandler(reg))
	r.Get("/metrics", handler.MetricsHandler(reg))
	r.Post("/service/add", handler.AddServiceHandler(reg))

	// 代理示例：/txstore/* → txstore
	txstoreSvc, _ := reg.GetService("txstore")
	r.Handle("/txstore/*", router.NewProxyHandler(txstoreSvc))

	log.Println("bridge starting on port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
