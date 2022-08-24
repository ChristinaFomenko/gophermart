package postgres

import (
	"context"
	"database/sql"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
	"go.uber.org/zap"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

type AuthPostgres struct {
	db  *sql.DB
	log *zap.Logger
}

func NewAuthPostgres(db *sql.DB, log *zap.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:  db,
		log: log,
	}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, user *model.User) (int, error) {
	stmt, err := a.db.PrepareContext(ctx,
		"INSERT INTO public.users(login, password) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result := stmt.QueryRowContext(ctx, user.Login, user.Password)
	var output sql.NullInt32
	_ = result.Scan(&output)
	if !output.Valid {
		return 0, errs.ConflictLoginError{
			Login: user.Login,
		}
	}
	userID := int(output.Int32)
	return userID, nil
}

func (a *AuthPostgres) GetUserID(ctx context.Context, user *model.User) (int, error) {
	row := a.db.QueryRowContext(ctx, "SELECT id FROM public.users WHERE login=$1 AND password=$2", user.Login, user.Password)
	var output sql.NullInt32
	_ = row.Scan(&output)
	if !output.Valid {
		return 0, errs.AuthenticationError{}
	}
	userID := int(output.Int32)
	return userID, nil
}
