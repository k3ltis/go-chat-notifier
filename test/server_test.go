package server_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	server "go-teams-notifier/internal/server"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Prepare post request body
	file, err := os.Open("data/backup_failure_warning.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		t.Error(err)
	}
	defer file.Close()
	body := io.Reader(file)

	// start server
	serverInstance := server.Run()
	time.Sleep(100 * time.Millisecond)

	// graceful shutdown server
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := serverInstance.Shutdown(ctx); err != nil {
			t.Fatalf("Server shutdown failed:%+v", err)
		}
	}()

	// test and assert
	resp, err := http.Post("http://localhost:8080/notification", "application/json", body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
