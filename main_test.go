package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

func Test_run(t *testing.T) {

	opts := options{
		Port: ":8081",
	}

	// Run the run function with the test options in a separate goroutine
	go func() {
		err := run(opts)
		assert.NoError(t, err)
	}()

	client := http.Client{Timeout: 3 * time.Second}
	file, err := os.Open("testdata/data.csv")
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
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

		url := fmt.Sprintf("http://localhost:8081/transactions")
		req, err := http.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() //nolint

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"status":"ok"}`+"\n", string(data))
	})

	t.Run("successful get", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8081/report")
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() //nolint
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, `{"grossRevenue":225,"expenses":72.93,"netRevenue":152.07}`+"\n", string(data))
	})

	// Wait for 5 seconds to let the run function run and then cancel it
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	signal.NotifyContext(ctx, os.Interrupt)
}
