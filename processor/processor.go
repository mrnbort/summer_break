package processor

import (
	"context"
	"fmt"
	"github.com/mrnbort/summer_break/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Proc allows to access transaction data
type Proc struct {
	mu           sync.RWMutex
	transactions []model.Transaction
}

// NewProc initiates and returns a slice for transaction data
func NewProc() *Proc {
	var transact []model.Transaction
	return &Proc{transactions: transact}
}

// ProcessTransactions adds new transactions to the transaction slice, thread-safe
func (p *Proc) ProcessTransactions(ctx context.Context, transactions []model.Transaction) error {
	// check ctx will be needed in case of non-memory (slow) storage
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	p.mu.Lock()
	p.transactions = append(p.transactions, transactions...)
	p.mu.Unlock()

	return nil
}

// ParseTransaction parses input csv record
func (p *Proc) ParseTransaction(rec []string) (model.Transaction, error) {
	amount, err := strconv.ParseFloat(strings.TrimSpace(rec[2]), 64)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("incorrect amount value %q: %w", rec[2], err)
	}

	transaction := model.Transaction{
		Type:   model.TrType(strings.TrimSpace(rec[1])), // TODO: add validation
		Amount: amount,
		Memo:   strings.TrimSpace(rec[3]),
	}
	transaction.Date, err = time.ParseInLocation("2006-01-02", rec[0], time.Local) // use the current location
	if err != nil {
		return model.Transaction{}, fmt.Errorf("invalid date %q: %w", rec[0], err)
	}
	return transaction, nil
}

// GenerateReport calculates revenue and expenses from transaction data and returns them
func (p *Proc) GenerateReport(ctx context.Context) (model.Report, error) {
	select {
	case <-ctx.Done():
		return model.Report{}, ctx.Err()
	default:
	}

	res := model.Report{}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, transaction := range p.transactions {
		switch transaction.Type {
		case model.Expense:
			res.Expenses += transaction.Amount
		case model.Income:
			res.GrossRevenue += transaction.Amount
		default:
			return model.Report{}, fmt.Errorf("unsupported transaction type %q", transaction.Type)
		}
	}
	res.NetRevenue = res.GrossRevenue - res.Expenses
	return res, nil
}
