package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	b, err := bot.New(os.Getenv("BOT_KEY"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "D4C is working...",
		})
	}
}
