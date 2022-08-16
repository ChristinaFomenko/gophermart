package storage

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
	"github.com/ChristinaFomenko/gophermart/internal/repository/psql"
)

type Transaction interface {
	BeginTx(context.Context) (*sql.Tx, error)
	AccrualOrderRepoContract
	WithdrawOrderRepoContract
}

type AuthRepoContract interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserID(ctx context.Context, user *model.User) (int, error)
}

type AccrualOrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number uint64) int
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type WithdrawOrderRepoContract interface {
	GetAccruals(tx *sql.Tx, ctx context.Context, UserID int) float32
	GetWithdrawals(tx *sql.Tx, ctx context.Context, UserID int) float32
	DeductPoints(tx *sql.Tx, ctx context.Context, order *model.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}

type Repository struct {
	Auth        AuthRepoContract
	Accrual     AccrualOrderRepoContract
	Withdraw    WithdrawOrderRepoContract
	Transaction Transaction
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{
		Auth:     postgres.NewAuthPostgres(db, log),
		Accrual:  postgres.NewAccrualOrderPostgres(db, log),
		Withdraw: postgres.NewWithdrawOrderPostgres(db, log),
	}
}
