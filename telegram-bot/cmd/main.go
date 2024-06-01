package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	teleBot, err := bot.New(os.Getenv("TELEGRAM_API_BOT"), opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	WebHookUrl := os.Getenv("WEBHOOK_URL")
	res, err := teleBot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: WebHookUrl,
	})

	if !res || err != nil {
		log.Fatalf("Error failed to set WebHook.")
	}

	go func() {
		err := http.ListenAndServe(":8080", teleBot.WebhookHandler())

		if err != nil {
			log.Fatalf("Failed to start server")
		}
	}()

	// Use StartWebhook instead of Start
	teleBot.StartWebhook(ctx)

	// call methods.DeleteWebhook if needed
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})

	if err != nil {
		log.Println("Failed to send reply to bot")
	}
}
