package main

import (
	"context"
	server "go-teams-notifier/internal/server"
)

func main() {
	server.Run(context.TODO())
}
