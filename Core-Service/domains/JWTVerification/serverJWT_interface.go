package JWTVerification

import (
	"context"

	"e-wallet-system-api/domains/JWTVerification/entities"
	"e-wallet-system-api/domains/JWTVerification/models/request"
)

type ServerJWTUseCase interface {
	Login(ctx context.Context, requests *request.LoginRequest) error
	CheckJWTServer(ctx context.Context, requests *request.CheckJWTRequest) (*entities.ServerJWT, error)
	Logout(ctx context.Context, requests *request.LogoutRequest) error
}

type ServerJWTRepository interface {
	FindServerJWTbyUserID(ctx context.Context, Id string) (*entities.ServerJWT, error)
	CreateServerJWT(ctx context.Context, accessToken string, userId string) (*entities.ServerJWT, error)
	UpdateServerJWT(ctx context.Context, ServerJWT *entities.ServerJWT) error
	DeleteServerJWT(ctx context.Context, userId string) error
}
