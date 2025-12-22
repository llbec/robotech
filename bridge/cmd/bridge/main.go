package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"bridge/internal/grpcx"
	"bridge/internal/health"
	"bridge/internal/httpx"
	"bridge/internal/metrics"
	"bridge/internal/registry"

	"google.golang.org/grpc/credentials"
)

func main() {
	reg := registry.New()

	metrics.Init()
	metrics.StartHTTP(":9100") // Prometheus metrics

	// Health check
	health.Start(reg, 5*time.Second)

	// HTTP server
	go func() {
		log.Println("HTTP bridge listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", httpx.NewHandler(reg)))
	}()

	// gRPC TLS server
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gRPC bridge (TLS) listening on :9090")
	log.Fatal(grpcx.NewServerWithTLS(reg, creds).Serve(lis))
}
