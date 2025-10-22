package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"robotech/httpapi"
	"robotech/logger"
)

func main() {
	// 定义命令行参数
	ip := flag.String("ip", "0.0.0.0", "HTTP server listen IP address")
	port := flag.Int("port", 8080, "HTTP server listen port")
	logfile := flag.String("log", "app.log", "Log file path")

	flag.Parse()

	// 初始化日志
	logger.Init(*logfile)

	// 初始化路由
	router := httpapi.NewRouter()

	addr := fmt.Sprintf("%s:%d", *ip, *port)
	log.Printf("🚀 HTTP server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
