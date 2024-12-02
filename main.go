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

var currentOptions = []bool{false, false, false}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithMessageTextHandler("/select", bot.MatchTypeExact, defaultHandler),
		bot.WithCallbackQueryDataHandler("btn_", bot.MatchTypePrefix, callbackHandler),
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
	go shorturl.Shorturl()
	b.Start(ctx)

}
func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	switch update.CallbackQuery.Data {
	case "btn_opt1":
		currentOptions[0] = !currentOptions[0]
	case "btn_opt2":
		currentOptions[1] = !currentOptions[1]
	case "btn_opt3":
		currentOptions[2] = !currentOptions[2]
	case "btn_select":
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   fmt.Sprintf("Selected options : %v", currentOptions),
		})
		return
	}
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: buildKeyboard(),
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
			Text:   "this is your shorten URL : http://localhost:4000/" + shortenURL,
		})
	}

}

func buildKeyboard() models.ReplyMarkup {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: buttonText("Opt 1", currentOptions[0]), CallbackData: "btn_opt1"},
				{Text: buttonText("Opt 2", currentOptions[1]), CallbackData: "btn_opt2"},
				{Text: buttonText("Opt 3", currentOptions[2]), CallbackData: "btn_opt3"},
			},
			{

				{Text: "Select", CallbackData: "btn_select"},
			},
		},
	}
	return kb
}
func buttonText(text string, opt bool) string {
	if opt {
		return "✅ " + text
	}
	return "❌ " + text
}
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "1 - /ShortURL for shorten URL Service 😡 \n 2- /FileConverter for Converting Files Service 🤬",
		ReplyMarkup: &models.ForceReply{
			ForceReply: true,
		},
	})
}
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " Select multiple optins for your command ",
		ReplyMarkup: buildKeyboard(),
	})
}
