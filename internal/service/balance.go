package service

import (
	"context"
	"database/sql"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

type WithdrawOrderRepoContract interface {
	GetAccruals(tx *sql.Tx, ctx context.Context, UserID int) float32
	GetWithdrawals(tx *sql.Tx, ctx context.Context, UserID int) float32
	DeductPoints(tx *sql.Tx, ctx context.Context, order *model.WithdrawOrder) error
	GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error)
}
type WithdrawOrderService struct {
	repo WithdrawOrderRepoContract
	log  *zap.Logger
}

func NewWithdrawOrderService(repo WithdrawOrderRepoContract, log *zap.Logger) *WithdrawOrderService {
	return &WithdrawOrderService{
		repo: repo,
		log:  log,
	}
}

func (w WithdrawOrderService) GetBalance(tx *sql.Tx, ctx context.Context, userID int) (float32, float32) {
	accruals := w.repo.GetAccruals(tx, ctx, userID)
	withdrawn := w.repo.GetWithdrawals(tx, ctx, userID)
	return accruals, withdrawn
}

func (w WithdrawOrderService) DeductionOfPoints(tx *sql.Tx, ctx context.Context, order *model.WithdrawOrder) error {
	accruals, withdrawn := w.GetBalance(tx, ctx, order.UserID)

	if order.Sum >= accruals-withdrawn {
		return errs.NotEnoughPoints{}
	}

	err := w.repo.DeductPoints(tx, ctx, order)
	if err != nil {
		w.log.Error("WithdrawOrderService.DeductionOfPoints: DeductPoints db error")
		return err
	}

	return nil
}

func (w *WithdrawOrderService) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error) {
	orders, err := w.repo.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		w.log.Error("WithdrawOrderService.GetWithdrawalOfPoints: GetWithdrawalOfPoints db error")
		return nil, err
	}
	return orders, nil
}
