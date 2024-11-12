package dto

type Balance struct {
	UserID        uint64
	CurrentAmount string
}

type UpdateBalance struct {
	UserID        uint64
	State         string
	Source        string
	Amount        string
	ParsedAmount  float64
	TransactionID string
}
