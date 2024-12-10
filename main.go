package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"telegram-bot/shorturl"

	FSM "telegram-bot/userState"

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
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, dynamicHandler)

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
	userID := update.Message.From.ID
	userState := FSM.GetUserFSM(userID)
	if userState.StateMachine.Is("idle") {
		userState.StateMachine.Event(ctx, "start")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a URL nigga",
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "We are already in the process of doing that shi ",
		})
	}
}

func dynamicHandler(ctx context.Context, b *bot.Bot, update *models.Update){
	userID := update.Message.From.ID
	userState := FSM.GetUserFSM(userID)
	if userState.StateMachine.Is("waiting_for_url") {
		url := update.Message.Text
		ShortenURL, err := shorturl.ShortenURLHandler(url)
		if err != nil {
			log.Println("error on the telegram side URL")
			return
		}
		userState.StateMachine.Event(ctx, "receive_url")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "this is your shorten URL : http://localhost:4000/" + ShortenURL,
		})
	}
	userState.StateMachine.Event(ctx, "rest")
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
	})
}
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " Select multiple optins for your command ",
		ReplyMarkup: buildKeyboard(),
	})
}
