package service

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
	"github.com/ChristinaFomenko/gophermart/internal/repository"
)

//type Transaction interface {
//	BeginTx(context.Context) (*sql.Tx, error)
//	AccrualOrderServiceContract
//	WithdrawOrderServiceContract
//}

type TxConnection interface {
	WithTx(context.Context, func(w WithdrawOrderServiceContract) error) error
}

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
	Auth         AuthServiceContract
	Accrual      AccrualOrderServiceContract
	Withdraw     WithdrawOrderServiceContract
	TxConnection TxConnection
}

func NewService(tx TxConnection, repo *storage.Repository, log *zap.Logger) *Service {
	return &Service{
		Auth:     NewAuthService(repo.Auth, log),
		Accrual:  NewAccrualOrderService(repo.Accrual, log),
		Withdraw: NewWithdrawOrderService(tx, repo.Withdraw, log),
	}
}
