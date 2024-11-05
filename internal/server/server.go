package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	api "go-teams-notifier/internal/api"
	discord "go-teams-notifier/internal/discord"
	openapi "go-teams-notifier/internal/generated/openapi/openapi"
)

const PORT = ":8080"

func Run(ctx context.Context) {
	// Setup notification receiver
	config := discord.NewDiscordClientConfig(os.Getenv("DISCORD_WEBHOOK_URL"))
	discordClient := discord.NewDiscordClient(config, discord.DefaultJSONMarshaller{}, discord.DefaultHTTPClient{})

	// setup service and router
	NotificationAPIService := api.NewNotificationAPIService(discordClient)
	NotificationAPIController := openapi.NewNotificationAPIController(NotificationAPIService)
	router := openapi.NewRouter(NotificationAPIController)

	// listen and serve
	fmt.Printf("Listen and serve on port %s\n", PORT)
	server := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	// Listen and serve in goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for server shutdown
	select {
	case <-ctx.Done():
		fmt.Println("Server shutdown via context done")
	case <-stopSignal():
		fmt.Println("Server shutdown via signal")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Server shutdown")
}

func stopSignal() <-chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	return stop
}
