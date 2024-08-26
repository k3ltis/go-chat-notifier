//go:build !integration
// +build !integration

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

func TestDiscordIntegration(t *testing.T) {
	// Prepare post request body
	file, err := os.Open("data/backup_failure_warning.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		t.Error(err)
	}
	defer file.Close()
	body := io.Reader(file)

	// Prepare context to shutdown server
	ctx, cancelServer := context.WithCancel(context.Background())
	defer cancelServer()

	// Start server
	go server.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	// test and assert
	resp, err := http.Post("http://localhost:8080/notification", "application/json", body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cancelServer()

	// Check for graceful exit
	select {
	case <-ctx.Done():
	default:
		t.Fatalf("Server did not shut down as expected")
	}
}
