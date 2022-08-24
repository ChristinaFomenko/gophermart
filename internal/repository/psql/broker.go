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

	stmt, err := tx.PrepareContext(ctx,
		"UPDATE public.accruals SET status=$1, amount=$2 WHERE order_num=$3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, order := range orderAccruals {
		_, err = stmt.ExecContext(ctx, order.Status, order.Accrual, order.Order)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
