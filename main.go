package main

import (
	bgrem "DirtyTaskDoneDirtFast/Background_Remover"
	Docx "DirtyTaskDoneDirtFast/Docconverter"
	fileconverter "DirtyTaskDoneDirtFast/FileConverter"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"github.com/mr-tron/base58"
)

// Define a map of commands to handler functions
var commandHandlers = map[string]func(ctx context.Context, b *bot.Bot, update *models.Update){
	"help":                     handleHelp,
	"we":                       handleSayHi,
	"nikmok":                   Wswoata3mok,
	"createwallet":             WalletCreation,
	"/transaction":             Transaction,
	"/checkbalance":            CheckBalance,
	"/imagebg":                 Images,
	"if you don't talk to her": talktuah,
	"/convert":                 Convert,
	"/convertpdf":              DocxConverter,
}

type ForTransaction struct {
	PrivateKey string
	SenderKey  string
}

var userfileType = make(map[int64]string)
var userState = make(map[int64]string)
var userInfo = make(map[int64]ForTransaction)

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

func Convert(ctx context.Context, b *bot.Bot, update *models.Update) {

	userId := update.Message.From.ID
	state, exist := userState[userId]
	if !exist {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Alright Send ur File here "})
		userState[userId] = "WaitingForFile"

	}
	switch state {

	case "WaitingForFile":
		if update.Message.Photo != nil {
			file := update.Message.Photo
			fileconverter.ConvertingFiles(file, b, ctx, update, "")
			delete(userState, userId)
		}
		if update.Message.Document != nil {
			file := update.Message.Document
			MimeType := file.MimeType
			fileconverter.ConvertingFiles(file, b, ctx, update, MimeType)
			delete(userState, userId)
		}

	}
}

func talktuah(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "  talk tuah    ",
	})

}
func Wswoata3mok(ctx context.Context, b *bot.Bot, update *models.Update) {

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   " MA GHADI NIKLK SOWA TA3 MOK A WELD 9A7BA    ",
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
		case "CheckingBalance":
			CheckBalance(ctx, b, update)

		case "WaitingForFile":
			Convert(ctx, b, update)
		case "waiting_for_doc":
			DocxConverter(ctx, b, update)

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
func DocxConverter(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	state, exist := userState[userId]
	if !exist {
		userState[userId] = "waiting_for_doc"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide the Docx",
		})
	}
	switch state {
	case "waiting_for_doc":
		if update.Message.Document != nil {
			// Check if the document is a .docx file
			if update.Message.Document.MimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
				// Proceed with document conversion
				Docx.ConvertingDocToPDF(update.Message.Document, ctx, b, update)

			} else {
				// If the document is not a .docx, inform the user
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Please send a .docx file.",
				})
			}
			delete(userState, userId)
		} else {
			// If no document is sent, inform the user
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Please send a .docx file.",
			})
			delete(userState, userId)
		}
	}

}
func CheckBalance(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	_, exists := userState[userId]

	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "PUT THE DAMMN PUBLIC KEY "})
		userState[userId] = "CheckingBalance"
		return
	}
	cmd := exec.Command("node", "./Solana_Stuff/Create_Public_Key.js", update.Message.Text)
	response, err := cmd.Output()
	if err != nil {
		log.Printf("error happend while checking the balance of a wallet %v ", err)
		delete(userState, userId)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: string(response)})
	delete(userState, userId)

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
		ValidationPublicKey, err := base58.Decode(update.Message.Text)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "error Happend make sure u typed correctly  ",
			})
			delete(userState, userId)
			delete(userInfo, userId)
			return
		}
		if len(ValidationPublicKey) != 64 {

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "PUT A VALID KEY NEXT TIME   ",
			})
			delete(userState, userId)
			delete(userInfo, userId)

			return
		}
		forTransaction.PrivateKey = update.Message.Text
		userState[userId] = "SenderKey"
		userInfo[userId] = forTransaction
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please Provide Sender key ",
		})
	case "SenderKey":
		forTransaction = userInfo[userId]
		ValidationPublicKey, err := base58.Decode(update.Message.Text)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "error Happend make sure u typed correctly  ",
			})
			return
		}
		if len(ValidationPublicKey) != 32 {

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "PUT A VALID KEY NEXT TIME   ",
			})
			delete(userState, userId)
			delete(userInfo, userId)

			return
		}
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
			delete(userState, userId)
			return
		}
		TransactioInformation := userInfo[userId]
		PrivateKey := TransactioInformation.PrivateKey
		fmt.Print("\n This is The Private kEy in the Main.go ", PrivateKey, "\n")
		SenderKey := TransactioInformation.SenderKey
		fmt.Print("\n This is The Sender kEy in the Main.go ", SenderKey, "\n")

		// TransactionCalling.CallingTransaction(ctx, b, update, PrivateKey, SenderKey, float32(floatval))
		Sol := strconv.FormatFloat(float64(floatval), 'f', 2, 32)

		cmd := exec.Command("node", "./Solana_Stuff/Transaction.js", PrivateKey, SenderKey, Sol)
		result, err := cmd.Output()
		if err != nil {
			fmt.Print("sra error f execCommand ", err)

		}
		signatureBase58 := string(result)

		fmt.Print("The result of the transaction is ", signatureBase58)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "heres ur signature ",
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   signatureBase58,
		})
		delete(userState, userId)

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
		log.Print("Error: Could not separate keys")
		return
	}
	publicKey := strings.TrimSpace(parts[0])
	privateKey := strings.TrimSpace(parts[1])
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: publicKey})
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: privateKey})

}
