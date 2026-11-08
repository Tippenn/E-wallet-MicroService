package usecase

import (
	"context"
	"errors"

	serverJWT "e-wallet-system-api/domains/JWTVerification"
	serverJWTRequest "e-wallet-system-api/domains/JWTVerification/models/request"
	"e-wallet-system-api/domains/users"
	"e-wallet-system-api/domains/users/entities"
	"e-wallet-system-api/domains/users/models/request"
	"e-wallet-system-api/domains/users/models/response"
	"e-wallet-system-api/shared/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo  users.UserRepository
	jwtUC serverJWT.ServerJWTUseCase
}

func NewUserUseCase(repo users.UserRepository, jwtUC serverJWT.ServerJWTUseCase) users.UserUseCase {
	return userUseCase{repo: repo, jwtUC: jwtUC}
}

// Login implements [users.UserUseCase].
func (u userUseCase) Login(ctx context.Context, requests *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := u.repo.FindUserbyEmail(requests.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requests.Password))
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	if user.Enabled == 0 {
		return nil, errors.New("user is disabled")
	}

	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		return nil, errors.New("token is not generated correctly")
	}

	loginRequest := serverJWTRequest.LoginRequest{
		UserId:      user.Id,
		AccessToken: token,
	}
	err = u.jwtUC.Login(ctx, &loginRequest)

	if err != nil {
		return nil, errors.New("token is not verified on the database")
	}

	return &response.LoginResponse{Token: token}, nil
}

// SignUp implements [users.UserUseCase].
func (u userUseCase) SignUp(ctx context.Context, requests *request.SignUpRequest) (*response.SignUpResponse, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(requests.Password), 10)
	if err != nil {
		return nil, err
	}

	uid := uuid.New().String()
	u.repo.CreateUser(ctx,
		uid,
		requests.Email,
		requests.Name,
		string(b),
	)

	return &response.SignUpResponse{
		UserId:  uid,
		Name:    requests.Name,
		Email:   requests.Email,
		Message: "User Registered Succesfully",
	}, nil
}

// FindUserbyUsername implements [users.UserUseCase].
func (u userUseCase) FindUserbyEmail(email string) (*entities.User, error) {
	user, err := u.repo.FindUserbyEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
