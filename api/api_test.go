package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mrnbort/summer_break/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestService_Run(t *testing.T) {
	done := make(chan struct{})
	go func() {
		<-done
		e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.NoError(t, e)
	}()
}

func TestService_handleTransactions(t *testing.T) {

	proc := &ProcessorMock{
		ProcessTransactionsFunc: func(ctx context.Context, trs []model.Transaction) error {
			return nil
		},
		ParseTransactionFunc: func(rec []string) (model.Transaction, error) {
			return model.Transaction{Amount: 123, Type: model.Income, Memo: "aaaa", Date: time.Now()}, nil
		},
	}

	svc := &Service{
		Processor: proc,
	}

	ts := httptest.NewServer(svc.routes())
	defer ts.Close()

	client := http.Client{Timeout: time.Second}
	file, err := os.Open("../testdata/data.csv")
	require.NoError(t, err)
	defer file.Close()

	t.Run("successful post", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileField, err := writer.CreateFormFile("file", "test.csv")
		require.NoError(t, err)
		_, err = io.Copy(fileField, file)
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)

		url := fmt.Sprintf("%s/transactions", ts.URL)
		req, err := http.NewRequest("POST", url, body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() //nolint

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"status":"ok"}`+"\n", string(data))
		require.Equal(t, 1, len(proc.ProcessTransactionsCalls()))
	})

	t.Run("failed post", func(t *testing.T) {
		proc.ProcessTransactionsFunc = func(ctx context.Context, trs []model.Transaction) error {
			return errors.New("oh oh")
		}

		_, err := file.Seek(0, io.SeekStart) // Reset the file cursor to the beginning
		require.NoError(t, err)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileField, err := writer.CreateFormFile("file", "test.csv")
		require.NoError(t, err)
		_, err = io.Copy(fileField, file)
		require.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		url := fmt.Sprintf("%s/transactions", ts.URL)
		req, err := http.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() //nolint
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"error":"oh oh"}`+"\n", string(data))
		require.Equal(t, 2, len(proc.ProcessTransactionsCalls()))
	})
}

func TestService_handleReport(t *testing.T) {
	proc := &ProcessorMock{
		GenerateReportFunc: func(ctx context.Context) (model.Report, error) {
			return model.Report{
				GrossRevenue: 20,
				Expenses:     30,
				NetRevenue:   40,
			}, nil
		},
	}

	svc := &Service{
		Processor: proc,
	}

	ts := httptest.NewServer(svc.routes())
	defer ts.Close()

	client := http.Client{Timeout: time.Second}

	t.Run("successful get", func(t *testing.T) {
		url := fmt.Sprintf("%s/report", ts.URL)
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"grossRevenue":20,"expenses":30,"netRevenue":40}`+"\n", string(data))
		require.Equal(t, 1, len(proc.GenerateReportCalls()))
	})

	t.Run("failed get", func(t *testing.T) {
		proc.GenerateReportFunc = func(ctx context.Context) (model.Report, error) {
			return model.Report{}, errors.New("oh oh")
		}
		url := fmt.Sprintf("%s/report", ts.URL)
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"error":"oh oh"}`+"\n", string(data))
		require.Equal(t, 2, len(proc.GenerateReportCalls()))
	})
}
