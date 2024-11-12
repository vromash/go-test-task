package handlers

import (
	"go_test_task/internal/dto"
	"go_test_task/internal/exchange"
	"go_test_task/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type BalanceHandler struct {
	service services.BalanceService
}

func NewBalanceHandler(balanceService services.BalanceService) BalanceHandler {
	return BalanceHandler{service: balanceService}
}

func (h BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &exchange.ErrorResponse{
			ErrorText: err.Error(),
		})
		return
	}

	balance, err := h.service.GetBalance(r.Context(), uint64(userID))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &exchange.ErrorResponse{
			ErrorText: err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapBalanceDTOToResponse(balance))
}

func (h BalanceHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	req := &exchange.UpdateBalanceRequest{}
	if err := render.Bind(r, req); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusOK)
}

func mapBalanceDTOToResponse(balance dto.Balance) *exchange.GetBalanceResponse {
	return &exchange.GetBalanceResponse{
		UserID:  balance.UserID,
		Balance: balance.Balance,
	}
}
