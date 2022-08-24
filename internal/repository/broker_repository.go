package storage

import (
	"context"
	"database/sql"
	psql "github.com/ChristinaFomenko/gophermart/internal/repository/psql"
	"go.uber.org/zap"

	modelBroker "github.com/ChristinaFomenko/gophermart/internal/broker/model"
)

type BrokerRepoContract interface {
	GetOrdersForProcessing(ctx context.Context, limit int) ([]modelBroker.Order, error)
	UpdateOrderAccruals(ctx context.Context, orderAccruals []modelBroker.OrderAccrual) error
}

type BrokerRepository struct {
	BrokerRepoContract
}

func NewBrokerRepository(db *sql.DB, log *zap.Logger) *BrokerRepository {
	return &BrokerRepository{
		BrokerRepoContract: psql.NewBrokerPostgres(db, log),
	}
}
