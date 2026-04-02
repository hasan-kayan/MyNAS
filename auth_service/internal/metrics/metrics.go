package metrics

import (
	"net/http"

	"github.com/Raos-Health/auth_service/internal/logger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Metrics struct {
	log logger.Logger
}

func New(log logger.Logger) *Metrics {
	return &Metrics{
		log: log,
	}
}

func (m *Metrics) RegisterGRPCMetrics(grpcServer *grpc.Server) {
	m.log.Info("Registering gRPC Prometheus metrics")
	grpc_prometheus.Register(grpcServer)
}

func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func (m *Metrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_prometheus.StreamServerInterceptor
}

func (m *Metrics) StartMetricsServer(address string) {
	m.log.Info("Starting Prometheus metrics server",
		zap.String("address", address),
	)

	go func() {
		if err := http.ListenAndServe(address, promhttp.Handler()); err != nil {
			m.log.Error("Metrics server failed",
				zap.Error(err),
			)
		}
	}()
}
