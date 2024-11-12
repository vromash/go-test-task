package handlers

import (
	"fmt"
	"go_test_task/internal/dto"
	"go_test_task/internal/exchange"
	"go_test_task/internal/models"
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
	userID, err := h.getUserID(r)
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	balance, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, mapBalanceDTOToResponse(balance))
}

func (h BalanceHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	source := r.Header.Get("Source-Type")
	if _, err := models.StringToTransactionSource(source); err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	req := &exchange.UpdateBalanceRequest{}
	if err := render.Bind(r, req); err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	if _, err := models.StringToTransactionState(req.State); err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	parsedAmount, err := h.parseUpdateBalanceAmount(req.Amount)
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	err = h.service.UpdateBalance(r.Context(), mapUpdateBalanceeRequestToDTO(*req, source, userID, parsedAmount))
	if err != nil {
		renderErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	render.JSON(w, r, &exchange.ErrorResponse{
		ErrorText: err.Error(),
	})
}

func mapBalanceDTOToResponse(balance dto.Balance) *exchange.GetBalanceResponse {
	return &exchange.GetBalanceResponse{
		UserID:  balance.UserID,
		Balance: balance.CurrentAmount,
	}
}

func mapUpdateBalanceeRequestToDTO(req exchange.UpdateBalanceRequest, source string, userID uint64, parsedAmount float64) dto.UpdateBalance {
	return dto.UpdateBalance{
		UserID:        userID,
		State:         req.State,
		Source:        source,
		Amount:        req.Amount,
		ParsedAmount:  parsedAmount,
		TransactionID: req.TransactionID,
	}
}

func (h BalanceHandler) getUserID(r *http.Request) (uint64, error) {
	userIDParam := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, err
	}

	if userID < 0 {
		return 0, fmt.Errorf("invalid user id")
	}

	return uint64(userID), nil
}

func (h BalanceHandler) parseUpdateBalanceAmount(amount string) (float64, error) {
	parsedAmount, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		return 0, err
	}

	if parsedAmount < 0 {
		return 0, fmt.Errorf("amount can't be negative")
	}

	return parsedAmount, nil
}
