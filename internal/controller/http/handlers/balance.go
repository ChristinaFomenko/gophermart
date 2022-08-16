package handler

import (
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

	tx, err := h.Service.Transaction.BeginTx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	accruals, withdraws := h.Service.Withdraw.GetBalance(tx, r.Context(), userID)

	b := balance{Current: accruals - withdraws, Withdrawn: withdraws}

	tx.Commit()

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

	tx, err := h.Service.Transaction.BeginTx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	err = h.Service.Withdraw.DeductionOfPoints(tx, r.Context(), order)

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusOK)
	case errs.NotEnoughPoints:
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	default:
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getWithdrawalOfPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.getUserIDFromToken(w, r, "handler.getCurrentBalance")
	if err != nil {
		return
	}

	orders, err := h.Service.Withdraw.GetWithdrawalOfPoints(r.Context(), userID)
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
