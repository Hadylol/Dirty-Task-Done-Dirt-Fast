package bgrem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ThisfunctiondoesSomething(photo interface{}, ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Printf("Type of passed value: %T\n", photo) // Prints the Go type
	fmt.Println("Reflect Type:", reflect.TypeOf(photo))
	switch f := photo.(type) {
	case *models.PhotoSize:
		print(f.FileID)
		file, err := b.GetFile(ctx, &bot.GetFileParams{
			FileID: f.FileID,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "PROCESSING THE FOKING IMAGE MAN WAIT "})

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
		localfile := fmt.Sprintf("../urshit%v.png", file.FileUniqueID)
		if _, err := Downloadimage(fileUrl, localfile); err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "makhdmtsh si zabi ma telecharjash 3adna ADSI"})

			return

		}
		cmd := exec.Command("python", "./Background_Remover/backgroundRemover.py", localfile, file.FileUniqueID)
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
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:   update.Message.Chat.ID,
			Document: &models.InputFileUpload{Filename: result, Data: bytes.NewReader(fileContent)},
		})

		os.Remove(localfile)
		os.Remove(result)

	case **models.Document:
		fmt.Print("  rani f document")

		file, err := b.GetFile(ctx, &bot.GetFileParams{
			FileID: (*f).FileID,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "PROCESSING THE FOKING IMAGE MAN WAIT "})

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
		if _, err := Downloadimage(fileUrl, localfile); err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "makhdmtsh si zabi ma telecharjash 3adna ADSI"})

			return

		}
		fmt.Print("Downloaded  the file")

		cmd := exec.Command("python", "./Background_Remover/backgroundRemover.py", localfile, file.FileUniqueID)
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

		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:   update.Message.Chat.ID,
			Document: &models.InputFileUpload{Filename: result, Data: bytes.NewReader(fileContent)},
		})

		os.Remove(localfile)
		os.Remove(result)

	default:
		fmt.Print(f)
	}

}
func Downloadimage(url, outputpath string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := make([]byte, 512)
	n, err := resp.Body.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Detect content type from the buffer
	contentType := http.DetectContentType(buf[:n])
	fmt.Printf("Detected content type from buffer: %s \n", contentType)

	// Reset the response body to copy the entire content
	bodyReader := io.MultiReader(bytes.NewReader(buf[:n]), resp.Body)
	Testingthis := strings.Split(outputpath, ".")
	Testingthecontentype := strings.Split(contentType, "/")
	FileType := Testingthis[1]

	if FileType == "png" {
		theoutputpath := fmt.Sprintf("%v.%v", Testingthis[0], Testingthecontentype[1])
		outputpath = theoutputpath
	}

	out, err := os.Create(outputpath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, bodyReader)
	return outputpath, nil

}
