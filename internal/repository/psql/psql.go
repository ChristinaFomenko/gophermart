package postgres

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Psql struct {
	DB *sql.DB
}

type ctxTxKey struct{}

func NewPsql(DSN string) (*Psql, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}
	pgx := &Psql{DB: db}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return pgx, nil
}

func (p *Psql) WithTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, alreadyHasTx := ctx.Value(ctxTxKey{}).(*sql.Tx)
	if !alreadyHasTx {
		tx, err = p.DB.Begin()
		if err != nil {

			return errors.WithStack(err)
		}
		ctx = context.WithValue(ctx, ctxTxKey{}, tx)
	}

	err = errors.WithStack(fn(ctx))

	if alreadyHasTx {

		return err
	}
	if err == nil {

		return errors.WithStack(tx.Commit())
	}

	tx.Rollback()

	return err
}

func (p *Psql) ExtractTx(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {

	return p.WithTx(ctx, func(ctx context.Context) error {
		tx := ctx.Value(ctxTxKey{}).(*sql.Tx)

		return errors.WithStack(fn(ctx, tx))
	})
}

func (p *Psql) Ping() error {
	if err := p.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (p *Psql) Init() error {
	_, err := p.DB.Exec(`CREATE TABLE IF NOT EXISTS users(
		    id SERIAL PRIMARY KEY,
    		login TEXT NOT NULL UNIQUE,
    		password TEXT NOT NULL,
    		"current" FLOAT NOT NULL DEFAULT 0,
        	withdrawal FLOAT NOT NULL DEFAULT 0
    		);

			CREATE TABLE IF NOT EXISTS orders(
				id BIGSERIAL PRIMARY KEY,
				order_num BIGINT UNIQUE,
				user_id INT NOT NULL,
				FOREIGN KEY (user_id) REFERENCES public.users (id));

			CREATE TABLE IF NOT EXISTS accruals(
				order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				status TEXT NOT NULL DEFAULT 'NEW',
				amount FLOAT DEFAULT 0,
				uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				FOREIGN KEY (user_id) REFERENCES public.users (id),
			    FOREIGN KEY (order_num) REFERENCES public.orders (order_num));

			CREATE TABLE IF NOT EXISTS withdrawals(
			    order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				amount FLOAT DEFAULT 0,
				processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				FOREIGN KEY (user_id) REFERENCES public.users (id),
			    FOREIGN KEY (order_num) REFERENCES public.orders (order_num));
`)

	if err != nil {
		return err
	}

	return nil
}
