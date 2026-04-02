package server

import (
	"net"

	authpb "github.com/Raos-Health/auth_service/api/auth"
	"github.com/Raos-Health/auth_service/internal/logger"
	"google.golang.org/grpc"
)

func RunGRPCServer(
	log logger.Logger,
	handler *AuthHandler,
) error {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server on :50051")

	grpcServer := grpc.NewServer()

	authpb.RegisterAuthenticationServiceServer(grpcServer, handler)

	return grpcServer.Serve(lis)
}
