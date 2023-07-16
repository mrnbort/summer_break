package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mrnbort/summer_break/model"
	"io"
	"log"
	"net/http"
	"time"
)

//go:generate moq -out processor_mock.go . Processor

// Service allows access to the transaction data
type Service struct {
	Processor    Processor
	Port         string
	httpServer   *http.Server
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

// Processor interface provides access to the functions that work with transaction data
type Processor interface {
	ParseTransaction(rec []string) (model.Transaction, error)
	ProcessTransactions(ctx context.Context, transactions []model.Transaction) error
	GenerateReport(ctx context.Context) (model.Report, error)
}

// JSON is a map alias, just for convenience
type JSON map[string]interface{}

// Run the listener and request's router, activates the rest server
func (s Service) Run(ctx context.Context) error {
	s.httpServer = &http.Server{
		Addr:         s.Port,
		Handler:      s.routes(),
		ReadTimeout:  s.ReadTimeOut,
		WriteTimeout: s.WriteTimeOut,
		IdleTimeout:  time.Second * 30, // don't allow hanging connection for long time
	}

	go func() {
		<-ctx.Done()
		log.Printf("[DEBUG] termination requested")
		err := s.httpServer.Close()
		if err != nil {
			log.Printf("[WARN] can't close server: %v", err)
		}
		log.Printf("[INFO] server closed")
	}()

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("service failed to run, err:%w", err)
	}
	return nil
}

func (s Service) routes() chi.Router {
	mux := chi.NewRouter()
	mux.Post("/transactions", s.handleTransactions)
	mux.Get("/report", s.handleReport)
	return mux
}

// POST /transactions
func (s Service) handleTransactions(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("[WARN] can't get file: %v", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, JSON{"error": err.Error()})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	transactions := []model.Transaction{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[WARN] failed to read line: %v", err)
			continue
		}

		// Skip lines with a different number of fields
		if len(record) < 4 {
			log.Printf("[WARN] skipping line with wrong number of fields: %v", record)
			continue
		}

		transaction, err := s.Processor.ParseTransaction(record)
		if err != nil {
			log.Printf("[WARN] failed to parse %v: %v", record, err)
			continue
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) <= 0 {
		log.Printf("[WARN] input file has no valid transations")
		render.Status(r, http.StatusBadRequest)
		return
	}

	err = s.Processor.ProcessTransactions(r.Context(), transactions)
	if err != nil {
		log.Printf("[WARN] can't process transactions: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, JSON{"error": err.Error()})
		return
	}

	render.JSON(w, r, JSON{"status": "ok"})
}

// GET /report
func (s Service) handleReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	report, err := s.Processor.GenerateReport(ctx)
	if err != nil {
		log.Printf("[WARN] can't generate report: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, JSON{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(report); err != nil {
		log.Printf("[WARN] can't encode report: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, JSON{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
