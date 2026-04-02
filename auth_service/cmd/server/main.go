package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	authpb "github.com/Raos-Health/auth_service/api/auth"
	"github.com/Raos-Health/auth_service/internal/config"
	"github.com/Raos-Health/auth_service/internal/logger"
	"github.com/Raos-Health/auth_service/internal/metrics"
	"github.com/Raos-Health/auth_service/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {

	// =============================
	// Initialize Logger
	// =============================
	logg, err := logger.New("development") // change to "production" in prod
	if err != nil {
		panic(err)
	}
	defer logg.Sync()

	// =============================
	// Bootstrap Log
	// =============================

	logg.Info("Bootstrapping Authentication Service",
		zap.String("service", "auth"),
		zap.String("version", "v1"),
	)

	// =============================
	// Load Configuration
	// =============================
	cfg, err := config.Load()
	if err != nil {
		logg.Error("Failed to load configuration")
		panic(err)
	}
	test := cfg.Server.GRPCPort
	logg.Info("Configuration loaded successfully", zap.String("GRPCPort To See Test", test))

	// =============================
	//  Create TCP Listener
	// =============================
	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		logg.Error("Failed to bind TCP listener",
			zap.Error(err),
		)
		os.Exit(1)
	}

	// =============================
	//  Create gRPC Server with Interceptors ((Metrics + Logging))
	// =============================
	metricsService := metrics.New(logg) // Create a new Metrics service instance with the logger
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metricsService.UnaryServerInterceptor()),
		grpc.StreamInterceptor(metricsService.StreamServerInterceptor()),
	)
	reflection.Register(grpcServer) // Enable gRPC reflection for easier debugging and tooling support
	metricsService.RegisterGRPCMetrics(grpcServer)
	metricsService.StartMetricsServer(":" + cfg.Server.MetricsPort)
	// =============================
	// Register Auth Handler
	// =============================
	authHandler := server.NewAuthHandler(logg)
	authpb.RegisterAuthenticationServiceServer(grpcServer, authHandler)

	// =============================
	//  Register gRPC Health Service
	// =============================
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// Mark service as SERVING
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// =============================
	// Start Server (Non-blocking)
	// =============================
	go func() {
		logg.Info("gRPC server started",
			zap.String("address", ":"+cfg.Server.GRPCPort),
		)
		if err := grpcServer.Serve(lis); err != nil {
			logg.Error("gRPC server failed",
				zap.Error(err),
			)
		}
	}()

	// =============================
	//  Graceful Shutdown Handling
	// =============================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // Wait for termination signal

	logg.Warn("Shutdown signal received")

	// Set health to NOT_SERVING
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	// Allow ongoing requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logg.Info("Server shutdown gracefully")
	case <-ctx.Done():
		logg.Warn("Graceful shutdown timed out, forcing stop")
		grpcServer.Stop()
	}

	logg.Info("Authentication Service stopped")
}
