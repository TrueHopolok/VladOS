// Analyze new and edited messages to create a pun with it.
//
// Takes message's suffix and tries to find a pun of a best match.
package pun

import (
	"log/slog"

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

// TODO
func handleNew(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "message", "pun/new")
	_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		Text:   "new reply",
		ChatID: update.Message.Chat.ChatID(),
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		},
	})
	return err
}

// TODO
func handleEdit(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "message", "pun/edit")
	_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		Text:   "edit reply",
		ChatID: update.EditedMessage.Chat.ChatID(),
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.EditedMessage.MessageID,
		},
	})
	return err
}
