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

// ValidType checks if the file type is supported
func ValidType(filetype string) bool {
	validTypes := map[string]bool{
		"GIF":  true,
		"JPEG": true,
		"PNG":  true,
		"BMP":  true,
		"TIFF": true,
		"DOC":  true,
		"DOCX": true,
		"RTF":  true,
		"ODT":  true,
		"PDF":  true,
	}
	return validTypes[filetype]
}

const sizeLimit = 5 * 1024 * 1024 //2MB
type ConversionHandler func(file *os.File, fileSize int64, buffer *bytes.Buffer, b *bot.Bot, ctx context.Context, update *models.Update)

// ConvertingFiles handles the file conversion based on type

// I need to Check if the Image is Ending with Png or Jpeg I just need to check the shit you know
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
		Button("GIF", []byte(fmt.Sprintf("%s | gif", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("JPEG", []byte(fmt.Sprintf("%s | jpeg", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("BMP", []byte(fmt.Sprintf("%s | bmp", fileID)), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("TIFF", []byte("tiff"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("DOC", []byte("doc"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("DOCX", []byte("docx"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("RTF", []byte("rtf"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Button("ODT", []byte("odt"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("PDF", []byte("pdf"), func(ctx context.Context, bot *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			processFile(b, ctx, update, data, MimeType)
		}).
		Row().
		Button("Cancel", []byte("cancel"), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Why waste my time man "})
		})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select the Type of File",
		ReplyMarkup: kb,
	})

}

// processFile handles downloading and processing the file
func processFile(b *bot.Bot, ctx context.Context, update *models.Update, fileID []byte, MimeType string) {
	theData := string(fileID)
	parts := strings.Split(theData, "|")
	fileId := parts[0]
	NewFileId := strings.TrimSpace(fileId)

	// extension := parts[1]
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
	convertedFile, err := os.Open(outputfile)
	if err != nil {
		log.Printf("The error happend at opening the FIle in FileConverter %v", err)
		return
	}
	defer convertedFile.Close()
	state, _ := convertedFile.Stat()
	//Switch Type to check the button Clicked and call the correct Function
	buffer := bytes.NewBuffer(make([]byte, 0, sizeLimit))

	switch parts[1] {
	case "gif":
		GifFuncConvert(localFile, convertedFile, state.Size(), buffer, b, ctx, update, file.FileUniqueID)

	case "jpeg":

	case "bmp":

	case "tiff":

	case "doc":

	case "DOCX":

	case "rtf":
	case "odt":
	case "pdf":

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

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing ffmpeg command: %v\n", err)
		fmt.Printf("Output: %s\n", output)
		return err
	}

	return nil

}
