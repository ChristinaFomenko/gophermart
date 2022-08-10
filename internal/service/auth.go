package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

const (
	secretKey = "be55d1079e6c6167118ac91318fe"
)

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserID(ctx context.Context, user *model.User) (int, error)
}

type AuthService struct {
	repo AuthRepoContract
	log  *zap.Logger
}

func NewAuthService(repo AuthRepoContract, log *zap.Logger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  log,
	}
}

func (auth *AuthService) CreateUser(ctx context.Context, user *model.User) error {
	user.Password = auth.generatePasswordHash(user.Password)
	userID, err := auth.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) AuthenticationUser(ctx context.Context, user *model.User) error {
	user.Password = auth.generatePasswordHash(user.Password)
	userID, err := auth.repo.GetUserID(ctx, user)
	if err != nil {
		return err
	}
	user.ID = userID
	return nil
}

func (auth *AuthService) GenerateToken(user *model.User, tokenAuth *jwtauth.JWTAuth) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
	return tokenString, err
}

func (auth *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}
