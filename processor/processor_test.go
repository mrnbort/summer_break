package processor

import (
	"context"
	"github.com/mrnbort/summer_break/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestProc_ParseTransactions(t *testing.T) {
	tests := []struct {
		name  string
		inp   []string
		out   model.Transaction
		isErr bool
	}{
		{"valid expense", []string{"2020-07-01", "Expense", "18.77", "Fuel"}, model.Transaction{
			Date:   time.Date(2020, 7, 1, 0, 0, 0, 0, time.Local),
			Memo:   "Fuel",
			Type:   model.Expense,
			Amount: 18.77,
		}, false},
		{"valid income", []string{"2020-07-04", "Income", "40.00", "347 Woodrow"}, model.Transaction{
			Date:   time.Date(2020, 7, 4, 0, 0, 0, 0, time.Local),
			Memo:   "347 Woodrow",
			Type:   model.Income,
			Amount: 40.00,
		}, false},
		{"another valid expense", []string{"2020-07-06", "Income", "35.00", "219 Pleasant"}, model.Transaction{
			Date:   time.Date(2020, 7, 6, 0, 0, 0, 0, time.Local),
			Memo:   "219 Pleasant",
			Type:   model.Income,
			Amount: 35.00,
		}, false},
		{"wrong day", []string{"2020-07-BAD", "Income", "35.00", "219 Pleasant"}, model.Transaction{}, true},
		{"wrong amount", []string{"2020-07-06", "Income", "xyz35.00", "219 Pleasant"}, model.Transaction{}, true},
		{"too little fields", []string{"2020-07-06", "Income", "xyz35.00"}, model.Transaction{}, true},
	}

	proc := Proc{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := proc.ParseTransaction(tt.inp)
			if tt.isErr {
				require.Error(t, err)
				return
			}
			assert.InDelta(t, out.Amount, tt.out.Amount, 0.0001)
			assert.Equal(t, out.Type, tt.out.Type)
			assert.Equal(t, out.Memo, tt.out.Memo)
			assert.Equal(t, out.Date, tt.out.Date)
		})
	}
}

func TestProc_ProcessTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	proc := &Proc{}

	err := proc.ProcessTransactions(ctx, []model.Transaction{
		{
			Date:   time.Date(2020, 7, 1, 0, 0, 0, 0, time.Local),
			Memo:   "Fuel",
			Type:   model.Expense,
			Amount: 18.77,
		},
		{
			Date:   time.Date(2020, 7, 4, 0, 0, 0, 0, time.Local),
			Memo:   "347 Woodrow",
			Type:   model.Income,
			Amount: 40.00,
		},
	})
	require.NoError(t, err)
	assert.Len(t, proc.transactions, 2)
}

func TestProc_GenerateReport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	transactions := []model.Transaction{
		{
			Date:   time.Date(2020, 7, 1, 0, 0, 0, 0, time.Local),
			Memo:   "Fuel",
			Type:   model.Expense,
			Amount: 18.77,
		},
		{
			Date:   time.Date(2020, 7, 4, 0, 0, 0, 0, time.Local),
			Memo:   "347 Woodrow",
			Type:   model.Income,
			Amount: 40.00,
		},
		{
			Date:   time.Date(2020, 7, 6, 0, 0, 0, 0, time.Local),
			Memo:   "219 Pleasant",
			Type:   model.Income,
			Amount: 35.00,
		},
	}

	proc := &Proc{
		transactions: transactions,
		mu:           sync.RWMutex{},
	}

	report, err := proc.GenerateReport(ctx)
	require.NoError(t, err)

	assert.InDelta(t, 75.00, report.GrossRevenue, 0.0001)
	assert.InDelta(t, 18.77, report.Expenses, 0.0001)
	assert.InDelta(t, 56.23, report.NetRevenue, 0.0001)
}
