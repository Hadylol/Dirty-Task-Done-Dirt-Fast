package fileconverter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"log"
	"os"

	bgrem "DirtyTaskDoneDirtFast/Background_Remover"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

const sizeLimit = 5 * 1024 * 1024 //5MB
type ConversionHandler func(file *os.File, fileSize int64, buffer *bytes.Buffer, b *bot.Bot, ctx context.Context, update *models.Update)

func ConvertingFiles(file interface{}, b *bot.Bot, ctx context.Context, update *models.Update, MimeType string) {

	var fileID string

	switch f := file.(type) {
	case []models.PhotoSize:
		fileID = f[len(f)-1].FileID
	case *models.Document:
		fileID = (*f).FileID
	default:
		fmt.Print("Still a problme going on here")

		return
	}
	kb := inline.New(b).
		Row().
		Button("JPEG", []byte(fmt.Sprintf("%s | jpeg", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("PDF", []byte(fmt.Sprintf("%s | pdf", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("BMP", []byte(fmt.Sprintf("%s | bmp", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("TIFF", []byte(fmt.Sprintf("%s | tiff", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("DOC", []byte(fmt.Sprintf("%s | doc", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("DOCX", []byte(fmt.Sprintf("%s | docx", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("RTF", []byte(fmt.Sprintf("%s | rtf", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("ODT", []byte(fmt.Sprintf("%s | odt", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("Cancel", []byte("cancel"), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Why waste my time man "})
		})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select the Type of File you want to convert to ",
		ReplyMarkup: kb,
	})

}

// processFile handles downloading and processing the file
func processFile(b *bot.Bot, ctx context.Context, update *models.Update, fileID []byte, MimeType string) {
	theData := string(fileID)
	parts := strings.Split(theData, "|")
	fileId := parts[0]

	NewFileId := strings.TrimSpace(fileId)
	fmt.Print("\n File Id : ", fileID)

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: NewFileId})
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error retrieving your file. Please try again.",
		})
		return
	}

	fileUrl := b.FileDownloadLink(&models.File{
		FileID:       file.FileID,
		FileSize:     file.FileSize,
		FileUniqueID: file.FileUniqueID,
		FilePath:     file.FilePath,
	})
	var localFile string
	if MimeType == "" {
		localfile := fmt.Sprintf("ConvertedFile_%v.jpeg", file.FileUniqueID)
		localFile = localfile
	} else {
		mimetype := strings.Split(MimeType, "/")
		thetype := mimetype[1]
		localfile := fmt.Sprintf("ConvertedFile_%v.%v", file.FileUniqueID, thetype)
		localFile = localfile
	}

	outputfile, err := bgrem.Downloadimage(fileUrl, localFile)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to download and process your file. Please try again.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("File successfully processed! Saved as %s", localFile),
	})
	fmt.Print(" \n TESTETEST \n")

	convertedFile, err := os.Open(outputfile)

	if err != nil {
		log.Printf("The error happend at opening the FIle in FileConverter %v", err)
		return
	}
	fmt.Print(" \n TESTETEST222 \n")

	defer convertedFile.Close()
	state, _ := convertedFile.Stat()
	//Switch Type to check the button Clicked and call the correct Function
	buffer := bytes.NewBuffer(make([]byte, 0, sizeLimit))
	fmt.Print(" \n TESTETEST4444 \n")
	fmt.Print(parts[1])

	switch strings.ToLower(strings.TrimSpace(parts[1])) {

	case "jpeg":
		GifFuncConvert(localFile, convertedFile, state.Size(), buffer, b, ctx, update, file.FileUniqueID)

	case "bmp":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "bmp",
		})
	case "tiff":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "tiff",
		})
	case "doc":
		fmt.Print(" \n alright the Switch Works here \n")

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "doc ",
		})
	case "docx":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "docx",
		})

	}

}

func GifFuncConvert(inputfile string, file *os.File, fileSize int64, buffer *bytes.Buffer, b *bot.Bot, ctx context.Context, update *models.Update, unique string) {
	file.Seek(0, io.SeekStart)

	// buf := make([]byte, 512)
	// if _, err := file.Read(buf); err != nil {
	// 	log.Printf("The file has been read my man : ")
	// 	return
	// }
	// file.Seek(0, io.SeekStart)
	// img, err := gif.Decode(file)

	// Name := fmt.Sprintf("SomethingIntheWay %v.gif", unique)
	// if err != nil {
	// 	log.Printf("Error in the GifFuncConver Function  %v", err)
	// 	return
	// }
	// fmt.Print("Checking here Its After the Name creation ")
	// outputfile, err := os.Create(Name)
	// if err != nil {
	// 	log.Printf("Error in the GifFuncConver Os.Create Function  %v", err)
	// 	return
	// }
	// fmt.Print("Checking here Its After  the outputfile of os.Create     ")

	// err = png.Encode(outputfile, img)
	// if err != nil {
	// 	log.Printf("Error in the GifFuncConver png.Encode Function  %v", err)
	// 	return

	// }
	theInput := strings.Split(inputfile, ".")
	outputfile := fmt.Sprintf("%vRANDOMSHITGO.png", theInput[0])

	CaptureImageFromAvideo(b, ctx, update, inputfile, outputfile)

	fmt.Println("\n GIF successfully converted to PNG!")

}

func CaptureImageFromAvideo(b *bot.Bot, ctx context.Context, update *models.Update, inputpath string, outputpath string) error {

	fmt.Print("YO YO YO YO YO \n", outputpath)
	fmt.Print("YO YO YO YO YO \n", inputpath)

	cmd := exec.Command("ffmpeg", "-i", inputpath, "-frames:v", "1", "-ss", "00:00:01.000", outputpath)

	fmt.Printf("Executing: %v\n", cmd.Args)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error in outputing the captureImageFromAvideo : %v\n", err)

	}

	fmt.Printf("Output: %s\n", output)

	return nil

}
