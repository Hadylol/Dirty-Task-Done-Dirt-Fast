package conveter

import (
	bgrem "DirtyTaskDoneDirtFast/Background_Remover"
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/net/html"
)

var (
	ErrTemplatePlaceholdersNotFound = errors.New("placeholders not found in template")
	ErrDelimitersNotPassed          = errors.New("delimiters did not passed")
	ErrBadTemplate                  = errors.New("can't open template file, it's broken")
	ErrPlaceholderWithWhitespaces   = errors.New("some placeholders template has leading or tailing whitespace")

	documentXmlPathInZip = "word/document.xml"
	xmlTextTag           = "t"
)

func getDocumentXmlReader(templateBytes []byte) (io.Reader, error) {
	templateReader := bytes.NewReader(templateBytes)
	zipReader, err := zip.NewReader(templateReader, int64(len(templateBytes)))
	if err != nil {
		return nil, ErrBadTemplate
	}
	file, err := zipReader.Open(documentXmlPathInZip)
	if err != nil {
		return nil, ErrBadTemplate
	}
	return file, nil
}

// Gets `word/document.xml` as string from given `docx` file (it's basically `zip`). Returns error
// if file is not a valid `docx`
func getAllXmlText(reader io.Reader) string {

	var output string
	tokenizer := html.NewTokenizer(reader)
	prevToken := tokenizer.Token()
loop:
	for {
		tok := tokenizer.Next()
		switch {
		case tok == html.ErrorToken:
			break loop // End of the document,  done
		case tok == html.StartTagToken:
			prevToken = tokenizer.Token()
		case tok == html.TextToken:
			if prevToken.Data == "script" {
				continue
			}
			TxtContent := html.UnescapeString(string(tokenizer.Text()))
			if len(TxtContent) > 0 {
				output += TxtContent
			}
		}
	}
	return output
}

const sizeLimit = 2 * 1024 * 1024 // 2 MB

func ConvertingDocToPDF(document *models.Document, ctx context.Context, b *bot.Bot, update *models.Update) {
	file, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: document.FileID,
	})
	if err != nil {
		log.Print("Yo err happend here ig in the DocxAndstuffConverter", err)

	}

	url := b.FileDownloadLink(&models.File{
		FileID:       file.FileID,
		FileSize:     file.FileSize,
		FileUniqueID: file.FileUniqueID,
		FilePath:     file.FilePath,
	})
	outputpath := fmt.Sprintf("TheDocAndStuffTypeofStuffhaha%v.docx", document.FileUniqueID)
	if _, err := bgrem.Downloadimage(url, outputpath); err != nil {
		log.Print("problem at downloading the doc man ngl ", err)

	}

	doc, err := os.Open(outputpath)
	if err != nil {
		log.Print("problem at opnening  the doc man ngl ", err)

	}
	defer doc.Close()
	docBytes, err := io.ReadAll(doc)
	if err != nil {
		log.Printf("problem reading the doc: %v", err)
		return
	}
	xmlReader, err := getDocumentXmlReader(docBytes)
	if err != nil {
		log.Printf("problem extracting document.xml: %v", err)
		return
	}
	extractedText := getAllXmlText(xmlReader)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, extractedText, "", "", false)
	outpdf := fmt.Sprintf("TheDocAndStuffTypeofStuffhaha%v.pdf", document.FileUniqueID)
	err = pdf.OutputFileAndClose(outpdf)
	if err != nil {
		log.Fatalf("Failed to save PDF file: %v", err)
	}

	fmt.Println("PDF file created successfully as ", outpdf)
	fileContent, errore := os.ReadFile("./" + outpdf)
	if errore != nil {
		fmt.Println("Error reading file:", errore)
		return
	}
	b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: &models.InputFileUpload{Filename: "SOMETHINGPDF.pdf", Data: bytes.NewReader(fileContent)},
	})
	os.Remove(outpdf)
	time.Sleep(15 * time.Second)

	errorr := os.Remove(outputpath)
	if errorr != nil {
		log.Printf("Error deleting file %s: %v", outputpath, err)
	}
}
