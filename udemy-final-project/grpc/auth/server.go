package auth

import (
	"context"
	"errors"
	"log"

	"github.com/gurkanindibay/udemy-rest-api/models"
	authpb "github.com/gurkanindibay/udemy-rest-api/proto/auth"
	"github.com/gurkanindibay/udemy-rest-api/utils"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := user.Save(); err != nil {
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

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Printf("Attempting login for user: %s", req.Email)

	verifiedUser, err := models.VerifyUserCredentials(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to verify user credentials: %v", err)
		return nil, err
	}
	if verifiedUser == nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(verifiedUser.Email, verifiedUser.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	return &authpb.LoginResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}
