package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

type TransactionState string

const (
	Win  TransactionState = "win"
	Lose TransactionState = "lose"
)

var stringToTransactionState = map[string]TransactionState{
	"win":  Win,
	"lose": Lose,
}

func (e *TransactionState) Scan(value interface{}) error {
	*e = TransactionState(value.([]byte))
	return nil
}

func (e TransactionState) Value() (driver.Value, error) {
	return string(e), nil
}

func StringToTransactionState(val string) (TransactionState, error) {
	if state, ok := stringToTransactionState[val]; ok {
		return state, nil
	}
	return "", fmt.Errorf("invalid transaction state value")
}

type TransactionSource string

const (
	Game    TransactionSource = "game"
	Server  TransactionSource = "server"
	Payment TransactionSource = "payment"
)

var stringToTransactionSource = map[string]TransactionSource{
	"game":    Game,
	"server":  Server,
	"payment": Payment,
}

func (e *TransactionSource) Scan(value interface{}) error {
	*e = TransactionSource(value.([]byte))
	return nil
}

func (e TransactionSource) Value() (driver.Value, error) {
	return string(e), nil
}

func StringToTransactionSource(val string) (TransactionSource, error) {
	if source, ok := stringToTransactionSource[val]; ok {
		return source, nil
	}
	return "", fmt.Errorf("invalid transaction source value")
}

type Transaction struct {
	gorm.Model
	UID    string
	UserID uint64
	Source TransactionSource `sql:"type:transaction_source"`
	TState TransactionState  `json:"t_state",sql:"type:transaction_state"`
	Amount string
}
