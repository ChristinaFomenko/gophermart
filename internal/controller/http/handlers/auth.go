package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/ChristinaFomenko/gophermart/internal/model"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
)

func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := h.readUserData(w, r, &user, "registration")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.CreateUser(ctx, &user)

	if errors.As(err, &errs.ConflictLoginError{}) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		h.log.Error("Handler.registration: CreateUser service error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.writeToken(w, &user, "registration")
}

func (h *Handler) authentication(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := h.readUserData(w, r, &user, "authentication")
	if err != nil {
		return
	}

	ctx := context.Background()
	err = h.Service.Auth.AuthenticationUser(ctx, &user)

	if errors.As(err, &errs.AuthenticationError{}) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		h.log.Error("Handler.authentication: AuthenticationUser service error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.writeToken(w, &user, "authentication")
}
