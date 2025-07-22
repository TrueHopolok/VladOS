package bot

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandGhoul Command = Command{
	InfoFull: `
 /ghoul
Starts from 1000, subtracts 7.
Result is outputed in the message. Then the process is repeated till the 0. 
 `,
	InfoBrief: "output 1000-7 loop",
	Handler: func(ctx *th.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "ghoul")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		_, _, args := tu.ParseCommand(update.Message.Text)
		if len(args) > 0 {
			_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, CmdInvalidArgsAmount()...))
			return err
		}
		for i := 1000; i > 7; i -= 7 {
			_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, tu.Entityf("%4d-7=%-3d", i, i-7).Blockquote()))
			if err != nil {
				return err
			}
		}
		return nil
	},
	Conversation: nil,
}
