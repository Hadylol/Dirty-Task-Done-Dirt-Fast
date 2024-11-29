package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"telegram-bot/shorturl"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
	}
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	b, err := bot.New(os.Getenv("BOT_KEY"), opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/ShortURL", bot.MatchTypeContains, URLHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	//	b.RegisterHandler(bot.HandlerTypeMessageText, "/help@bosukeTest_bot", bot.MatchTypeExact, helpHandler)

	b.Start(ctx)
}
func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})

}
func URLHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	fakeURL := update.Message.Text
	fmt.Println(fakeURL)
	URL := strings.Fields(fakeURL)[1]
	log.Println("this is the URl provided", URL)
	if URL == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "please Provide a URL !",
		})
	} else {
		shortenURL, err := shorturl.ShortenURLHandler(URL)
		if err != nil {
			log.Println("failed to insert the data", err)
		}
		log.Println("the shorten URL is ", shortenURL)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "this is your shorten URL : " + shortenURL,
		})
	}

}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "button_1"},
				{Text: "Button 2", CallbackData: "button_2"},
			},
			{

				{Text: "Button 3", CallbackData: "button_3"},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " /help for all commands 😃👍",
		ReplyMarkup: kb,
	})
}
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "1 - /ShortURL for shorten URL Service 😡 \n 2- /FileConverter for Converting Files Service 🤬",
	})
}
