package grpcx

import (
	"io"
	"log"
	"strings"

	"bridge/internal/metrics"
	"bridge/internal/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewServer(reg *registry.Registry) *grpc.Server {
	return grpc.NewServer(
		grpc.UnknownServiceHandler(proxyHandler(reg)),
	)
}

func NewServerWithTLS(reg *registry.Registry, creds credentials.TransportCredentials) *grpc.Server {
	return grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnknownServiceHandler(proxyHandler(reg)),
	)
}

func proxyHandler(reg *registry.Registry) grpc.StreamHandler {
	return func(srv any, serverStream grpc.ServerStream) error {
		fullMethod, ok := grpc.MethodFromServerStream(serverStream)
		if !ok {
			return status.Errorf(status.Code(grpc.ErrServerStopped), "unknown method")
		}
		serviceName := parseServiceName(fullMethod)
		metrics.Requests.WithLabelValues(serviceName).Inc()
		log.Println("[bridge] gRPC call:", fullMethod)

		target := reg.Pick(serviceName)
		if target == nil {
			return status.Errorf(status.Code(grpc.ErrServerStopped), "no backend")
		}

		if !target.Limiter.Allow() {
			metrics.RateLimit.WithLabelValues(serviceName).Inc()
			return status.Errorf(status.Code(grpc.ErrServerStopped), "rate limit exceeded")
		}

		conn, err := grpc.NewClient(target.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return status.Errorf(status.Code(grpc.ErrServerStopped), "downstream connect failed: %v", err)
		}
		defer conn.Close()

		ctx := serverStream.Context()
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		stream, err := grpc.NewClientStream(ctx, &_StreamDesc, conn, fullMethod)
		if err != nil {
			return status.Errorf(status.Code(grpc.ErrServerStopped), "create downstream stream failed: %v", err)
		}

		errCh := make(chan error, 2)

		// client -> downstream
		go func() {
			for {
				req := new([]byte)
				if err := serverStream.RecvMsg(req); err != nil {
					errCh <- err
					return
				}
				if err := stream.SendMsg(req); err != nil {
					errCh <- err
					return
				}
			}
		}()

		// downstream -> client
		go func() {
			for {
				resp := new([]byte)
				if err := stream.RecvMsg(resp); err != nil {
					errCh <- err
					return
				}
				if err := serverStream.SendMsg(resp); err != nil {
					errCh <- err
					return
				}
			}
		}()

		for i := 0; i < 2; i++ {
			err := <-errCh
			if err == io.EOF {
				continue
			} else if err != nil {
				return err
			}
		}
		return nil
	}
}

func parseServiceName(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	parts := strings.Split(fullMethod, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

var _StreamDesc = grpc.StreamDesc{
	StreamName:    "Unknown",
	ServerStreams: true,
	ClientStreams: true,
}
