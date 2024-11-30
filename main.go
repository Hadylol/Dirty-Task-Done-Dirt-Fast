package main

import (
	bgrem "DirtyTaskDoneDirtFast/Background_Remover"
	TransactionCalling "DirtyTaskDoneDirtFast/Solana_Stuff"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

// Define a map of commands to handler functions
var commandHandlers = map[string]func(ctx context.Context, b *bot.Bot, update *models.Update){
	"help":         handleHelp,
	"we":           handleSayHi,
	"Toyo":         Toyo,
	"nikmok":       Wswoata3mok,
	"createWallet": WalletCreation,
	"/Transaction": Transaction,

	"/ImageBg": Images,
}

type ForTransaction struct {
	PrivateKey string
	SenderKey  string
}

var userState = make(map[int64]string)
var userInfo = make(map[int64]ForTransaction)

func UpdateUserState(id int64, state string) {

	userState[id] = state

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
func Toyo(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "saheb li khdemni w na9sh",
	})

}
func Wswoata3mok(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "sowa Ta3 MOK",
	})

}
func dynamichandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID

	state, exists := userState[userId]
	if exists {
		switch state {
		case "waiting_for_Images":
			Images(ctx, b, update)

		case "Transaction":
			Transaction(ctx, b, update)
		case "SenderKey":
			Transaction(ctx, b, update)
		case "Sol":
			Transaction(ctx, b, update)
		}

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
func Transaction(ctx context.Context, b *bot.Bot, update *models.Update) {
	var forTransaction ForTransaction
	userId := update.Message.From.ID
	state, exists := userState[userId]
	if !exists {
		userState[userId] = "Transaction"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide Private key ",
		})
		return
	}
	switch state {
	case "Transaction":
		forTransaction.PrivateKey = update.Message.Text
		userState[userId] = "SenderKey"
		userInfo[userId] = forTransaction
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide Sender key ",
		})
	case "SenderKey":
		forTransaction = userInfo[userId]
		forTransaction.SenderKey = update.Message.Text
		userState[userId] = "Sol"
		userInfo[userId] = forTransaction
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide How much solana you want to send ",
		})
	case "Sol":
		floatval, err := strconv.ParseFloat(update.Message.Text, 32)
		if err != nil {
			fmt.Print("\n", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Please put a number man  ",
			})
			return
		}
		TransactioInformation := userInfo[userId]
		PrivateKey := TransactioInformation.PrivateKey
		fmt.Print("\n This is The Private kEy in the Main.go ", PrivateKey, "\n")
		SenderKey := TransactioInformation.SenderKey
		fmt.Print("\n This is The Sender kEy in the Main.go ", SenderKey, "\n")

		TransactionCalling.CallingTransaction(ctx, b, update, PrivateKey, SenderKey, float32(floatval))
		delete(userInfo, userId)

	}
}
func Images(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	state, exists := userState[userId]
	if !exists {
		userState[userId] = "waiting_for_Images"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide the Image",
		})

		return
	}
	switch state {
	case "waiting_for_Images":
		fmt.Print(update.Message.From.ID, "\n")
		if update.Message.Photo != nil {
			photo := update.Message.Photo[len(update.Message.Photo)-1]
			go bgrem.ThisfunctiondoesSomething(&photo, ctx, b, update)

		}
		if update.Message.Document != nil && (update.Message.Document.MimeType == "image/jpeg" || update.Message.Document.MimeType == "image/png") {
			photo := update.Message.Document
			go bgrem.ThisfunctiondoesSomething(&photo, ctx, b, update)

		}
		delete(userState, userId)

	}
}
func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Available commands:\n/help - Show this message\n/we - Greet the bot!",
	})
}

func handleSayHi(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "we mrbk",
	})
}
func WalletCreation(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := exec.Command("node", "./Solana_Stuff/Create_Public_Key.js")
	output, err := cmd.Output()
	if err != nil {
		fmt.Print("u wanna break the glass ceieling hillary i sense it ", err)
	}
	data := string(output)
	parts := strings.Split(data, "\n")
	if len(parts) < 2 {
		fmt.Println("Error: Could not separate keys")
		return
	}
	publicKey := strings.TrimSpace(parts[0])
	privateKey := strings.TrimSpace(parts[1])
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: publicKey})
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: privateKey})

}
