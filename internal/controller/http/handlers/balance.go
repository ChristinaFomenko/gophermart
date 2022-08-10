package handler

import (
	"context"
	"encoding/json"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
	"io"
	"net/http"

	"github.com/ChristinaFomenko/gophermart/internal/model"
)

type balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

//getCurrentBalance GET /api/user/balance - получение текущего баланса пользователя
func (h *Handler) getCurrentBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.getUserIDFromToken(w, r, "handler.getCurrentBalance")
	if err != nil {
		return
	}

	ctx := context.Background()
	accruals, withdraws := h.Service.Withdraw.GetBalance(ctx, userID)

	b := balance{Current: accruals - withdraws, Withdrawn: withdraws}

	output, err := json.Marshal(b)
	if err != nil {
		h.log.Error("Handler.getCurrentBalance: json write error")
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

//deductionOfPoints POST /api/user/balance/withdraw - запрос на списание средств
func (h *Handler) deductionOfPoints(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromToken(w, r, "handler.deductionOfPoints")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Handler.deductionOfPoints: body read error")
		http.Error(w, "wrong input data", http.StatusInternalServerError)
		return
	}

	var order *model.WithdrawOrder
	err = json.Unmarshal(body, &order)
	if err != nil {
		h.log.Error("Handler.deductionOfPoints: json read error")
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}

	order.UserID = userID

	ctx := context.Background()
	err = h.Service.Withdraw.DeductionOfPoints(ctx, order)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusOK)
	case errs.NotEnoughPoints:
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	default:
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
	}
}

func (h *Handler) getWithdrawalOfPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.getUserIDFromToken(w, r, "handler.getCurrentBalance")
	if err != nil {
		return
	}

	ctx := context.Background()
	orders, err := h.Service.Withdraw.GetWithdrawalOfPoints(ctx, userID)
	if err != nil {
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	output, err := json.Marshal(orders)
	if err != nil {
		h.log.Error("Handler.getWithdrawalOfPoints: json marshal error")
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
