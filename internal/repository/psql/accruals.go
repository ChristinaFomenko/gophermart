package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

type AccrualOrderPostgres struct {
	db  *sql.DB
	log *zap.Logger
}

func NewAccrualOrderPostgres(db *sql.DB, log *zap.Logger) *AccrualOrderPostgres {
	return &AccrualOrderPostgres{
		db:  db,
		log: log,
	}
}

func (a *AccrualOrderPostgres) SaveOrder(ctx context.Context, order *model.AccrualOrder) (err error) {
	order.UploadedAt = time.Now()

	var userID int
	var current float32

	checkCurrent, err := a.db.PrepareContext(ctx, "SELECT id, current FROM users WHERE id = $1;")
	if err != nil {
		return err
	}

	if err = checkCurrent.QueryRowContext(ctx, order.UserID).Scan(&userID, &current); err != nil {
		return err
	}

	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			txError := tx.Rollback()
			if txError != nil {
				err = fmt.Errorf("accruals SaveOrder rollback error %s: %s", txError.Error(), err.Error())
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders(order_num,user_id) VALUES ($1,$2)", order.Number, order.UserID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO accruals(order_num,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4)",
		order.Number, order.UserID, order.Status.String(), order.UploadedAt)
	if err != nil {
		return err
	}

	if order.Status == model.StatusPROCESSED {
		current = current + order.Accrual
		updateAccrual, err := a.db.PrepareContext(ctx, "UPDATE users SET current = $1 WHERE id = $2;")
		if err != nil {
			return err
		}

		txUpdateAccrual := tx.StmtContext(ctx, updateAccrual)

		_, err = txUpdateAccrual.ExecContext(ctx, current, userID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (a *AccrualOrderPostgres) GetUserIDByNumberOrder(ctx context.Context, number uint64) int {
	row := a.db.QueryRowContext(ctx, "SELECT user_id FROM public.accruals WHERE order_num=$1", number)
	var userID int
	_ = row.Scan(&userID)

	return userID
}

func (a *AccrualOrderPostgres) GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error) {
	rows, err := a.db.QueryContext(ctx, "SELECT order_num, status, amount, uploaded_at FROM public.accruals WHERE user_id =$1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.AccrualOrder
	for rows.Next() {
		var order model.AccrualOrder
		var status string
		err = rows.Scan(&order.Number, &status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		order.Status, err = model.GetStatus(status)
		if err != nil {
			a.log.Error("broker db GetUploadedOrders")
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
