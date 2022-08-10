package service

import (
	"context"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

type AccrualOrderRepoContract interface {
	SaveOrder(ctx context.Context, order *model.AccrualOrder) error
	GetUserIDByNumberOrder(ctx context.Context, number uint64) int
	GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error)
}

type AccrualOrderService struct {
	repo AccrualOrderRepoContract
	log  *zap.Logger
}

func NewAccrualOrderService(repo AccrualOrderRepoContract, log *zap.Logger) *AccrualOrderService {
	return &AccrualOrderService{
		repo: repo,
		log:  log,
	}
}

func (a *AccrualOrderService) LoadOrder(ctx context.Context, numOrder uint64, userID int) error {

	if !a.Check(numOrder) {
		return errs.CheckError{}
	}

	order := model.AccrualOrder{
		Number: numOrder,
		UserID: userID,
		Status: model.StatusNEW,
	}

	userIDinDB := a.repo.GetUserIDByNumberOrder(ctx, order.Number)
	if userIDinDB != 0 {
		if userIDinDB == order.UserID {
			return errs.OrderAlreadyUploadedCurrentUserError{}
		} else {
			return errs.OrderAlreadyUploadedAnotherUserError{}
		}
	}

	err := a.repo.SaveOrder(ctx, &order)
	if err != nil {
		a.log.Error("AccrualOrderService.LoadOrder: SaveOrder db error")
		return err
	}

	return nil
}

func (a *AccrualOrderService) Check(number uint64) bool {
	var sum uint64

	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 {
			sum += cur
			number = number / 10
			continue
		}
		cur = cur * 2
		if cur > 9 {
			cur = cur - 9
		}
		sum += cur
		number = number / 10
	}

	return sum%10 == 0
}

func (a *AccrualOrderService) GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error) {
	orders, err := a.repo.GetUploadedOrders(ctx, userID)
	if err != nil {
		a.log.Error("AccrualOrderService.GetUploadedOrders: GetUploadedOrders db error")
		return nil, err
	}
	return orders, nil
}
