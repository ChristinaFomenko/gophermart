package handler

import (
	"encoding/json"
	errs "github.com/ChristinaFomenko/gophermart/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) loadOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromToken(w, r, "handler.loadOrders")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Handler.loadOrders: body read error")
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		h.log.Info("Handler.loadOrders: body empty")
		http.Error(w, "incorrect input data", http.StatusBadRequest)
		return
	}
	strBody := string(body)
	numOrder, err := strconv.ParseUint(strBody, 0, 64)
	if err != nil {
		h.log.Error("Handler.loadOrders: ParseUint number order error")
		http.Error(w, "wrong input data", http.StatusBadRequest)
		return
	}

	//statusCode := http.StatusAccepted

	err = h.Service.Accrual.LoadOrder(r.Context(), numOrder, userID)

	//if err != nil {
	//	if !errors.Is(err, errs.ErrOrderAlreadyUploadedCurrentUser) {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	statusCode = http.StatusOK
	//
	//	if !errors.Is(err, errs.ErrOrderAlreadyUploadedAnotherUser) {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	statusCode = http.StatusConflict
	//	if !errors.Is(err, errs.ErrInvalidOrderNumberForm) {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	statusCode = http.StatusUnprocessableEntity
	//
	//}

	switch err.(type) {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	case errs.OrderAlreadyUploadedCurrentUserError:
		http.Error(w, err.Error(), http.StatusOK)
		return
	case errs.OrderAlreadyUploadedAnotherUserError:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case errs.CheckError:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	default:
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
	}

}

func (h *Handler) getUploadedOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	userID, err := h.getUserIDFromToken(w, r, "handler.getUploadedOrders")
	if err != nil {
		return
	}

	orders, err := h.Service.Accrual.GetUploadedOrders(r.Context(), userID)
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
		h.log.Error("Handler.getUploadedOrders: json marshal error")
		http.Error(w, errs.InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write(output)
}
