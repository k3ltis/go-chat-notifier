package discord

type DiscordClientConfig struct {
	WebhookUrl string
}

func NewDiscordClientConfig(webhookUrl string) *DiscordClientConfig {
	return &DiscordClientConfig{WebhookUrl: webhookUrl}
}
