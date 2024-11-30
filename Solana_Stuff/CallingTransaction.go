package TransactionCalling

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ForTransaction struct {
	PrivateKey string
	SenderKey  string
}

func CallingTransaction(ctx context.Context, b *bot.Bot, update *models.Update, Privatekey string, SenderKey string, sol float32) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: Privatekey})
	if err != nil {
		fmt.Print(err)

	}
	fmt.Print(Privatekey)

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: SenderKey})
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "wee"})
	return
}
