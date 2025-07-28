// Contain a telegram bot logic to use on the server.
//
// Basicly: BRAIN of the VladOS.
package bot

import (
	"context"
	"log/slog"
	"os"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

var handler *th.BotHandler

// Connects all connectors and middlewares into given [github.com/mymmrac/telego/telegohandler.BotHandler] to serve.
//   - [LoggerMiddleware],
//   - [ConnectConversation],
//   - [ConnectCommands].
func ConnectAll(bh *th.BotHandler) {
	bh.Use(LoggerMiddleware)
	ConnectConversation(bh)
	ConnectCommands(bh)
	ConnectJokes(bh)
}

// Provides small bot handler middleware to connect for logs purposes using [log/slog] package.
func LoggerMiddleware(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "status", "START")
	defer slog.Debug("bot handler", "upd", update.UpdateID, "status", "FINISH")
	err := ctx.Next(update)
	if err != nil {
		slog.Error("bot handler", "upd", update.UpdateID, "status", "FAILED", "error", err)
	}
	return err
}

// Initialize a bot and starts it with handlers connected via [ConnectAll].
// Additionaly clears all DB dynamic tables (those that requires restarting every launch).
//
// Will stop execution of a previous bot in case it was working previously.
func Start(botErrorChan chan error) error {
	if err := Stop(); err != nil {
		return err
	}

	if err := dbconvo.Clear(); err != nil {
		return err
	}

	rawToken, err := os.ReadFile(cfg.Get().BotTokenPath)
	if err != nil {
		return err
	}

	bot, err := telego.NewBot(string(rawToken), telego.WithDiscardLogger())
	if err != nil {
		return err
	}

	ctx := context.Background()
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}

	handler, err = th.NewBotHandler(bot, updates)
	if err != nil {
		return err
	}

	ConnectAll(handler)

	go func() {
		botErrorChan <- handler.Start()
	}()
	select {
	case err := <-botErrorChan:
		return err
	default:
		return nil
	}
}

// Stop package's global bot from receiving and handling any updates.
func Stop() error {
	if handler == nil {
		return nil
	}
	return handler.Stop()
}
