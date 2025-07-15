// Contain a telegram bot logic to use on the server.
//
// Basicly: BRAIN of the VladOS.
package bot

import (
	"context"
	"os"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

var botUpdHandler *th.BotHandler

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// TODO
func ConnectAll(bot *telego.Bot) error {
	return nil
}

// Initialize for a whole pacakge, connecting all handlers via [ConnectAll].
// Return an error on initialization or after any repeated function call.
func Start(botErrorChan chan error) error {
	if err := Stop(); err != nil {
		return err
	}

	rawToken, err := os.ReadFile(cfg.Get().BotTokenPath)
	if err != nil {
		return err
	}

	bot, err := telego.NewBot(string(rawToken), telego.WithDefaultDebugLogger())
	if err != nil {
		return err
	}

	ctx := context.Background()
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return nil
	}

	botUpdHandler, err = th.NewBotHandler(bot, updates)
	return err
}

// Stop bot from receiving and handling the updates.
func Stop() error {
	if botUpdHandler == nil {
		return nil
	}
	return botUpdHandler.Stop()
}
