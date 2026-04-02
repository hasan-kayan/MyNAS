package server

import (
	"context"

	authpb "github.com/Raos-Health/auth_service/api/auth"
	"github.com/Raos-Health/auth_service/internal/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authpb.UnimplementedAuthenticationServiceServer
	log logger.Logger
}

func NewAuthHandler(log logger.Logger) *AuthHandler {
	return &AuthHandler{
		log: log,
	}
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {

	h.log.Info("Login request received",
		zap.String("email", req.Email),
	)

	return &authpb.LoginResponse{}, nil
}
