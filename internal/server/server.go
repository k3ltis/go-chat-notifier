package server

import (
	"fmt"
	"net/http"
	"os"

	api "go-teams-notifier/internal/api"
	discord "go-teams-notifier/internal/discord"
	openapi "go-teams-notifier/internal/generated/openapi/openapi"
)

const PORT = ":8080"

func Run() *http.Server {
	// Setup notification receiver
	config := discord.NewDiscordClientConfig(os.Getenv("DISCORD_WEBHOOK_URL"))
	discordClient := discord.NewDiscordClient(config)

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
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}
