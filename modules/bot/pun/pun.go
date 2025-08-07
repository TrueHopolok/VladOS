// Analyze new and edited messages to create a pun with it.
//
// Takes message's suffix and tries to find a pun of a best match.
package pun

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Connect a handler that analyze the message and find a joke / pun for the suffix of that message.
func ConnectJokes(bh *th.BotHandler) {
	// slog.Debug("bot handler", "upd", update.UpdateID, "message", "joke")
	ph := bh.Group(th.AnyMessage(), th.AnyEditedMessage())
	ph.Handle(handleNew, th.AnyMessage())
	ph.Handle(handleEdit, th.AnyEditedMessage())
}

// TODO
func handleNew(ctx *th.Context, update telego.Update) error {
	return nil
}

// TODO
func handleEdit(ctx *th.Context, update telego.Update) error {
	return nil
}
