package service

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
	"github.com/ChristinaFomenko/gophermart/internal/repository"
)

type AuthServiceContract interface {
	CreateUser(ctx context.Context, user *model.User) error
	AuthenticationUser(ctx context.Context, user *model.User) error
	GenerateToken(user *model.User, tokenAuth *jwtauth.JWTAuth) (string, error)
}

type AccrualOrderServiceContract interface {
	LoadOrder(ctx context.Context, numOrder uint64, userID int) error
	Check(number uint64) bool
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type WithdrawOrderServiceContract interface {
	DeductionOfPoints(ctx context.Context, order *model.WithdrawOrder) error
	GetBalance(ctx context.Context, userID int) (float32, float32)
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}

type Service struct {
	Auth     AuthServiceContract
	Accrual  AccrualOrderServiceContract
	Withdraw WithdrawOrderServiceContract
}

func NewService(repo *storage.Repository, log *zap.Logger) *Service {
	return &Service{
		Auth:     NewAuthService(repo.Auth, log),
		Accrual:  NewAccrualOrderService(repo.Accrual, log),
		Withdraw: NewWithdrawOrderService(repo.Withdraw, log),
	}
}
