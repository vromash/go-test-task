package exchange

import "net/http"

type GetBalanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
}

type UpdateBalanceRequest struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transactionId"`
}

func (u *UpdateBalanceRequest) Bind(r *http.Request) error {
	return nil
}

func (u *UpdateBalanceRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UpdateBalanceResponse struct {
}

type ErrorResponse struct {
	ErrorText string `json:"errorText"`
}
