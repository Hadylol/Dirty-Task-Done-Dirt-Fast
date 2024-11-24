package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

// Define a map of commands to handler functions
var commandHandlers = map[string]func(ctx context.Context, b *bot.Bot, update *models.Update){
	"/help":  handleHelp,
	"/sayHi": handleSayHi,
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
	// b.RegisterHandlerMatchFunc(matchfunc , )
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Photo != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Processessing the fucking image just wait cunt ",
		})
		photo := update.Message.Photo[len(update.Message.Photo)-1]

		thisfunctiondoesSomething(&photo, ctx, b, update)

	}
	if update.Message != nil && update.Message.Text != "" {
		if handler, exists := commandHandlers[update.Message.Text]; exists {
			handler(ctx, b, update) // Call the appropriate handler
			return
		}
	}

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

func thisfunctiondoesSomething(photo *models.PhotoSize, ctx context.Context, b *bot.Bot, update *models.Update) {
	file, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: photo.FileID,
	})
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "can't ma man "})
		return

	}
	fileUrl := b.FileDownloadLink(&models.File{
		FileID:       file.FileID,
		FileSize:     file.FileSize,
		FileUniqueID: file.FileUniqueID,
		FilePath:     file.FilePath,
	})
	localfile := fmt.Sprintf("urshit%v.png", file.FileUniqueID)
	if err := downloadimage(fileUrl, localfile); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "makhdmtsh si zabi ma telecharjash 3adna ADSI"})

		return

	}
	cmd := exec.Command("python", "backgroundRemover.py", localfile, file.FileUniqueID)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Python script error:", err)
		return
	}

	var builder strings.Builder
	for _, val := range output {
		builder.WriteByte(byte(val))
	}

	// Trim and clean up the result
	result := strings.TrimSpace(builder.String())
	fmt.Printf("Resulting filename: [%s]\n", result)

	// Verify if the file exists
	if _, err := os.Stat("./" + result); os.IsNotExist(err) {
		fmt.Printf("File does not exist: %s\n", result)
		return
	}

	// Read the file content
	fileContent, errore := os.ReadFile("./" + result)
	if errore != nil {
		fmt.Println("Error reading file:", errore)
		return
	}

	// Send the photo
	message, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileUpload{Filename: result, Data: bytes.NewReader(fileContent)},
		Caption: "New backgroundless photo ",
	})
	if err != nil {
		fmt.Println("Error sending photo:", err)
		return
	}
	os.Remove(result)
	os.Remove(localfile)

	fmt.Println("Photo sent successfully:", message)
}
func downloadimage(url, outputpath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(outputpath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err

}
