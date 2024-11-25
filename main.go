package main

import (
	bgrem "DirtyTaskDoneDirtFast/Background_Remover"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

// Define a map of commands to handler functions
var commandHandlers = map[string]func(ctx context.Context, b *bot.Bot, update *models.Update){
	"help": handleHelp,
	"we":   handleSayHi,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, _ := bot.New(os.Getenv("BOT_KEY"), opts...)
	b.RegisterHandlerMatchFunc(matchfunc, dynamichandler)
	b.Start(ctx)
}

func matchfunc(update *models.Update) bool {
	if update.Message == nil {
		return false

	}
	return update.Message.Text != "" || len(update.Message.Photo) > 0 || update.Message.Document != nil

}
func dynamichandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Photo != nil {
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		bgrem.ThisfunctiondoesSomething(&photo, ctx, b, update)

	}
	if update.Message != nil && update.Message.Text != "" {
		if handler, exists := commandHandlers[update.Message.Text]; exists {
			handler(ctx, b, update) // Call the appropriate handler
			return
		}
	}

}
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Print("buying poor socks i'll create jobs tearing down ur mom ")

}

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Available commands:\n/help - Show this message\n/sayHi - Greet the bot!",
	})
}

func handleSayHi(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "we mrbk",
	})
}
