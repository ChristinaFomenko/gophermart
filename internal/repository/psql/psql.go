package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Psql struct {
	DB *sql.DB
}

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

func (p *Psql) Ping() error {
	if err := p.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (p *Psql) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return p.DB.BeginTx(ctx, nil)
}

func (p *Psql) Init() error {
	_, err := p.DB.Exec(`CREATE TABLE IF NOT EXISTS public.users(
		    id SERIAL PRIMARY KEY,
    		login TEXT NOT NULL UNIQUE,
    		password TEXT NOT NULL
            current NUMERIC(10, 2) DEFAULT 0,
			withdrawn NUMERIC(10, 2) DEFAULT 0);

			CREATE TABLE IF NOT EXISTS public.orders(
				id BIGSERIAL PRIMARY KEY,
				order_num BIGINT UNIQUE,
				user_id INT NOT NULL,
				FOREIGN KEY (user_id) REFERENCES public.users (id));

			CREATE TABLE IF NOT EXISTS public.accruals(
				order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				status TEXT NOT NULL DEFAULT 'NEW',
				amount REAL DEFAULT 0,
				uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				FOREIGN KEY (user_id) REFERENCES public.users (id),
			    FOREIGN KEY (order_num) REFERENCES public.orders (order_num));

			CREATE TABLE IF NOT EXISTS withdrawals(
			    order_num BIGINT PRIMARY KEY,
				user_id INT NOT NULL,
				amount REAL DEFAULT 0,
				processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				FOREIGN KEY (user_id) REFERENCES public.users (id),
			    FOREIGN KEY (order_num) REFERENCES public.orders (order_num));
`)

	if err != nil {
		return err
	}

	return nil
}
