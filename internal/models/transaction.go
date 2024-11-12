package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type TransactionState string

const (
	Win  TransactionState = "win"
	Lose TransactionState = "lost"
)

func (e *TransactionState) Scan(value interface{}) error {
	*e = TransactionState(value.([]byte))
	return nil
}

func (e TransactionState) Value() (driver.Value, error) {
	return string(e), nil
}

type TransactionSource string

const (
	Game    TransactionSource = "game"
	Server  TransactionSource = "server"
	Payment TransactionSource = "payment"
)

func (e *TransactionSource) Scan(value interface{}) error {
	*e = TransactionSource(value.([]byte))
	return nil
}

func (e TransactionSource) Value() (driver.Value, error) {
	return string(e), nil
}

type Transaction struct {
	gorm.Model
	UserID uint64
	Source TransactionSource `sql:"type:transaction_source"`
	TState TransactionState  `json:"t_state",sql:"type:transaction_state"`
	Amount string
}
