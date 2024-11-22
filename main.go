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
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(thisfunctiondoesSomething),
	}
	bot, err := bot.New(os.Getenv("BOT_KEY"), opts...)
	if err != nil {
		panic(err)

	}
	bot.Start(ctx)
}
func thisfunctiondoesSomething(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "I'M RACIST",
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,

		Text: "🙊",
	})

}
