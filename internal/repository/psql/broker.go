package postgres

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/broker/model"
)

type BrokerPostgres struct {
	db  *sql.DB
	log *zap.Logger
}

func NewBrokerPostgres(db *sql.DB, log *zap.Logger) *BrokerPostgres {
	return &BrokerPostgres{
		db:  db,
		log: log,
	}
}

func (b *BrokerPostgres) GetOrdersForProcessing(ctx context.Context, limit int) ([]model.Order, error) {
	rows, err := b.db.QueryContext(ctx, "SELECT order_num, status FROM public.accruals WHERE status=$1 OR status=$2 LIMIT $3", model.StatusNEW.String(), model.StatusPROCESSING.String(), limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		var status string
		err = rows.Scan(&order.Number, &status)
		if err != nil {
			return nil, err
		}
		order.Status, err = model.GetStatus(status)
		if err != nil {
			b.log.Error("broker db GetOrdersForProcessing")
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (b *BrokerPostgres) UpdateOrderAccruals(ctx context.Context, orderAccruals []model.OrderAccrual) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	updateAccrualStmt, err := tx.PrepareContext(ctx,
		"UPDATE accruals SET status=$1, amount=$2 WHERE order_num=$3")
	if err != nil {
		return err
	}
	defer updateAccrualStmt.Close()

	updateUserStmt, err := b.db.PrepareContext(ctx, "UPDATE users SET current = $1 WHERE id =$2;")
	if err != nil {
		return err
	}

	txUpdateUserStmt := tx.StmtContext(ctx, updateUserStmt)

	for _, order := range orderAccruals {
		_, err = updateAccrualStmt.ExecContext(ctx, order.Status, order.Accrual, order.Order)
		if err != nil {
			return err
		}

		if order.Status == model.StatusPROCESSED {
			var current float32
			if err = tx.QueryRowContext(ctx, "SELECT current FROM users WHERE id = $1;", order.UserID).Scan(&current); err != nil {
				return err
			}
			current = current + order.Accrual

			_, err = txUpdateUserStmt.ExecContext(ctx, current, order.UserID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
