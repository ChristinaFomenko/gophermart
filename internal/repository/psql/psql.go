package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Psql struct {
	DB *sql.DB
}

type contextKey int

const txKey contextKey = iota

func NewPsql(DSN string) (*Psql, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}
	pgxd := &Psql{DB: db}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return pgxd, nil
}

func (p *Psql) Ping() error {
	if err := p.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (p *Psql) Conn(ctx context.Context) pgxtype.Querier {
	if tx := extractTx(ctx); tx != nil {
		return tx
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

func (p *Psql) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	if maybeTx := extractTx(ctx); maybeTx != nil {
		return txFunc(ctx)
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if errRollback := tx.Rollback(); errRollback != nil {
			if !errors.Is(errRollback, pgx.ErrTxClosed) {
			}
			return
		}
	}()

	err = txFunc(injectTx(ctx, tx))
	if err != nil {
		return err
	}

	// if no error, commit
	if errCommit := tx.Commit(); errCommit != nil {
		return errCommit
	}
	return nil
}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return nil
}
