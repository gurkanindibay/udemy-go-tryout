// Package auth provides gRPC authentication service implementation.
package auth

import (
	"context"
	"errors"
	"log"

	authpb "github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/proto/auth"
	"github.com/gurkanindibay/udemy-go-tryout/udemy-final-project/services"
)

// Server implements the gRPC AuthService server
type Server struct {
	authpb.UnimplementedAuthServiceServer
	userService services.UserService
	authService services.AuthService
}

// NewAuthServer creates a new AuthServer instance
func NewAuthServer(userService services.UserService, authService services.AuthService) *Server {
	return &Server{
		userService: userService,
		authService: authService,
	}
}

// Register handles user registration via gRPC
func (s *Server) Register(_ context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user, err := s.userService.Register(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return nil, err
	}

	return &authpb.RegisterResponse{
		User: &authpb.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

// Login handles user authentication via gRPC
func (s *Server) Login(_ context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Printf("Attempting login for user: %s", req.Email)

	verifiedUser, err := s.userService.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to verify user credentials: %v", err)
		return nil, err
	}
	if verifiedUser == nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.authService.GenerateToken(verifiedUser.Email, verifiedUser.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	return &authpb.LoginResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}
