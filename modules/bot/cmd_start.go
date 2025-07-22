package bot

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStart(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "start")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	_, _, args := tu.ParseCommand(update.Message.Text)
	if len(args) > 0 {
		_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, CmdInvalidArgsAmount()...))
		return err
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID,
		tu.Entity("Hello, "), tu.Entity("user.\n").Bold(),
		tu.Entity("I am bot "), tu.Entity("VladOS.\nVlad Operation System.\n").Bold(),
		tu.Entity(`
I am a project that combines:
 - Telegram bot;
 - Reincarnaction of the AllEgg bot from Discord;
 - Webpage to control and view bot activities.

Bot will react to any message and will try to find a pun for suffix of the message.

Also has additional functional via commands.
Type /help for more info about them.
`)))
	return err
}
