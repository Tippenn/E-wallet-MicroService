package usecase

import (
	"context"
	"errors"
	"time"

	"e-wallet-system-api/domains/JWTVerification"
	"e-wallet-system-api/domains/JWTVerification/entities"
	"e-wallet-system-api/domains/JWTVerification/models/request"
)

type serverJWTUseCase struct {
	repo JWTVerification.ServerJWTRepository
}

func NewServerJWTUseCase(repo JWTVerification.ServerJWTRepository) JWTVerification.ServerJWTUseCase {
	return serverJWTUseCase{repo: repo}
}

// Login implements [JWTVerification.ServerJWTUseCase].
func (s serverJWTUseCase) Login(ctx context.Context, requests *request.LoginRequest) error {
	_, err := s.repo.CreateServerJWT(ctx, requests.AccessToken, requests.UserId)

	if err != nil {
		return err
	}

	return nil
}

// Logout implements [JWTVerification.ServerJWTUseCase].
func (s serverJWTUseCase) Logout(ctx context.Context, requests *request.LogoutRequest) error {
	jwtData, err := s.repo.FindServerJWTbyUserID(ctx, requests.UserId)
	if err != nil {
		return err
	}

	err = s.repo.DeleteServerJWT(ctx, jwtData.UserId)

	if err != nil {
		return err
	}

	return nil
}

// CheckJWTServer implements [JWTVerification.ServerJWTUseCase].
func (s serverJWTUseCase) CheckJWTServer(ctx context.Context, requests *request.CheckJWTRequest) (*entities.ServerJWT, error) {
	jwtData, err := s.repo.FindServerJWTbyUserID(ctx, requests.UserId)
	if err != nil {
		return nil, err
	}

	if jwtData.Revoked {
		return nil, errors.New("token Revoked")
	}

	if time.Now().After(jwtData.ExpiresAt) {
		return nil, errors.New("token Expired")
	}

	return jwtData, nil
}
