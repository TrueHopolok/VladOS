// Analyze new and edited messages to create a pun with it.
//
// Takes message's suffix and tries to find a pun of a best match.
package pun

import (
	"log/slog"

	"github.com/TrueHopolok/VladOS/modules/db/dbpun"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Connect a handler that analyze the message and find a joke / pun for the suffix of that message.
func ConnectJokes(bh *th.BotHandler) {
	ph := bh.Group(th.Or(th.AnyMessage(), th.AnyEditedMessage()))
	ph.Handle(handleNew, th.AnyMessage())
	ph.Handle(handleEdit, th.AnyEditedMessage())
}

func handleNew(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "message", "pun/new")
	pun, err := dbpun.Answer(update.Message.Text)
	if err != nil {
		return err
	}
	if pun != "" {
		_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			Text:   pun,
			ChatID: update.Message.Chat.ChatID(),
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	}
	return err
}

func handleEdit(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "message", "pun/edit")
	pun, err := dbpun.Answer(update.Message.Text)
	if err != nil {
		return err
	}
	if pun != "" {
		_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			Text:   pun,
			ChatID: update.EditedMessage.Chat.ChatID(),
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.EditedMessage.MessageID,
			},
		})
	}
	return err
}
