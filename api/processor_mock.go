// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package api

import (
	"context"
	"github.com/mrnbort/summer_break/model"
	"sync"
)

// Ensure, that ProcessorMock does implement Processor.
// If this is not the case, regenerate this file with moq.
var _ Processor = &ProcessorMock{}

// ProcessorMock is a mock implementation of Processor.
//
//	func TestSomethingThatUsesProcessor(t *testing.T) {
//
//		// make and configure a mocked Processor
//		mockedProcessor := &ProcessorMock{
//			GenerateReportFunc: func(ctx context.Context) (model.Report, error) {
//				panic("mock out the GenerateReport method")
//			},
//			ParseTransactionFunc: func(rec []string) (model.Transaction, error) {
//				panic("mock out the ParseTransaction method")
//			},
//			ProcessTransactionsFunc: func(ctx context.Context, transactions []model.Transaction) error {
//				panic("mock out the ProcessTransactions method")
//			},
//		}
//
//		// use mockedProcessor in code that requires Processor
//		// and then make assertions.
//
//	}
type ProcessorMock struct {
	// GenerateReportFunc mocks the GenerateReport method.
	GenerateReportFunc func(ctx context.Context) (model.Report, error)

	// ParseTransactionFunc mocks the ParseTransaction method.
	ParseTransactionFunc func(rec []string) (model.Transaction, error)

	// ProcessTransactionsFunc mocks the ProcessTransactions method.
	ProcessTransactionsFunc func(ctx context.Context, transactions []model.Transaction) error

	// calls tracks calls to the methods.
	calls struct {
		// GenerateReport holds details about calls to the GenerateReport method.
		GenerateReport []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// ParseTransaction holds details about calls to the ParseTransaction method.
		ParseTransaction []struct {
			// Rec is the rec argument value.
			Rec []string
		}
		// ProcessTransactions holds details about calls to the ProcessTransactions method.
		ProcessTransactions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Transactions is the transactions argument value.
			Transactions []model.Transaction
		}
	}
	lockGenerateReport      sync.RWMutex
	lockParseTransaction    sync.RWMutex
	lockProcessTransactions sync.RWMutex
}

// GenerateReport calls GenerateReportFunc.
func (mock *ProcessorMock) GenerateReport(ctx context.Context) (model.Report, error) {
	if mock.GenerateReportFunc == nil {
		panic("ProcessorMock.GenerateReportFunc: method is nil but Processor.GenerateReport was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGenerateReport.Lock()
	mock.calls.GenerateReport = append(mock.calls.GenerateReport, callInfo)
	mock.lockGenerateReport.Unlock()
	return mock.GenerateReportFunc(ctx)
}

// GenerateReportCalls gets all the calls that were made to GenerateReport.
// Check the length with:
//
//	len(mockedProcessor.GenerateReportCalls())
func (mock *ProcessorMock) GenerateReportCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGenerateReport.RLock()
	calls = mock.calls.GenerateReport
	mock.lockGenerateReport.RUnlock()
	return calls
}

// ParseTransaction calls ParseTransactionFunc.
func (mock *ProcessorMock) ParseTransaction(rec []string) (model.Transaction, error) {
	if mock.ParseTransactionFunc == nil {
		panic("ProcessorMock.ParseTransactionFunc: method is nil but Processor.ParseTransaction was just called")
	}
	callInfo := struct {
		Rec []string
	}{
		Rec: rec,
	}
	mock.lockParseTransaction.Lock()
	mock.calls.ParseTransaction = append(mock.calls.ParseTransaction, callInfo)
	mock.lockParseTransaction.Unlock()
	return mock.ParseTransactionFunc(rec)
}

// ParseTransactionCalls gets all the calls that were made to ParseTransaction.
// Check the length with:
//
//	len(mockedProcessor.ParseTransactionCalls())
func (mock *ProcessorMock) ParseTransactionCalls() []struct {
	Rec []string
} {
	var calls []struct {
		Rec []string
	}
	mock.lockParseTransaction.RLock()
	calls = mock.calls.ParseTransaction
	mock.lockParseTransaction.RUnlock()
	return calls
}

// ProcessTransactions calls ProcessTransactionsFunc.
func (mock *ProcessorMock) ProcessTransactions(ctx context.Context, transactions []model.Transaction) error {
	if mock.ProcessTransactionsFunc == nil {
		panic("ProcessorMock.ProcessTransactionsFunc: method is nil but Processor.ProcessTransactions was just called")
	}
	callInfo := struct {
		Ctx          context.Context
		Transactions []model.Transaction
	}{
		Ctx:          ctx,
		Transactions: transactions,
	}
	mock.lockProcessTransactions.Lock()
	mock.calls.ProcessTransactions = append(mock.calls.ProcessTransactions, callInfo)
	mock.lockProcessTransactions.Unlock()
	return mock.ProcessTransactionsFunc(ctx, transactions)
}

// ProcessTransactionsCalls gets all the calls that were made to ProcessTransactions.
// Check the length with:
//
//	len(mockedProcessor.ProcessTransactionsCalls())
func (mock *ProcessorMock) ProcessTransactionsCalls() []struct {
	Ctx          context.Context
	Transactions []model.Transaction
} {
	var calls []struct {
		Ctx          context.Context
		Transactions []model.Transaction
	}
	mock.lockProcessTransactions.RLock()
	calls = mock.calls.ProcessTransactions
	mock.lockProcessTransactions.RUnlock()
	return calls
}
