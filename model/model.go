package model

import "time"

// TrType represents transaction type
type TrType string

// enum of all transaction types
const (
	Expense TrType = TrType("Expense")
	Income  TrType = TrType("Income")
)

// Transaction creates a transaction to save
type Transaction struct {
	Date   time.Time
	Type   TrType
	Amount float64
	Memo   string
}

// Report with revenue and expenses to return to user
type Report struct {
	GrossRevenue float64 `json:"grossRevenue"`
	Expenses     float64 `json:"expenses"`
	NetRevenue   float64 `json:"netRevenue"`
}
