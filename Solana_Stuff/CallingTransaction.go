package TransactionCalling

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CallingTransaction(ctx context.Context, b *bot.Bot, update *models.Update) {
	VariableGetter := update.Message.Text

}
