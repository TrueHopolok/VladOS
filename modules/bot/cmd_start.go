package bot

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStart(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handle", "command", "help")
	bot := ctx.Bot()
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(update.Message.Chat.ChatID(),
		tu.Entity("Hello, "), tu.Entity("user.\n").Bold(),
		tu.Entity("I am bot "), tu.Entity("VladOS.\nVlad Operation System.\n").Bold(),
		tu.Entity(`
I am a project that combines:
 - Telegram bot;
 - Reincarnaction of the AllEgg bot from Discord;
 - Webpage to control and view bot activities.

Type /help for more info about the functional.
`)))
	return err
}
