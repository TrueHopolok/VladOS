package cmd

import (
	"log/slog"
	"strconv"

	"github.com/TrueHopolok/VladOS/modules/db/dbtip"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandTip Command = Command{
	InfoBrief: "output some advice/tip/quote",
	InfoFull: `
 /tip
Output a random advice, tip or quote.

 /tip <advice_id>
Output a tip with id equal to given one. 
`,
	handler: func(ctx *th.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "tip")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		_, _, cmdArgs := tu.ParseCommand(update.Message.Text)
		var msgText []tu.MessageEntityCollection
		switch len(cmdArgs) {
		case 0:
			tip_text, author, tip_id, err := dbtip.Rand()
			if err != nil {
				return err
			}
			msgText = getTipText(tip_id, tip_text, author)
		case 1:
			tip_id, err := strconv.Atoi(cmdArgs[0])
			if err != nil || tip_id < 0 {
				_, err := bot.SendMessage(ctx, tu.Message(chatID, "Given argument is invalid, please enter the valid non-negative number.\nFor more info type:\n /help tip\n /help"))
				return err
			}
			tip_text, author, found, err := dbtip.Get(tip_id)
			if err != nil {
				return err
			}
			if !found {
				_, err := bot.SendMessage(ctx, tu.Message(chatID, "The advice/tip with given id does not exists."))
				return err
			}
			msgText = getTipText(tip_id, tip_text, author)
		default:
			_, err := bot.SendMessage(ctx, tu.Message(chatID, "Too many arguments are given for the command.\nFor more info type:\n /help tip\n /help"))
			return err
		}
		_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
		return err
	},
	conversation: nil,
}

func getTipText(tip_id int, tip_text, author string) []tu.MessageEntityCollection {
	return []tu.MessageEntityCollection{
		tu.Entityf("TIP #%d\n", tip_id).Bold(),
		tu.Entity(tip_text).Blockquote(),
		tu.Entityf("\n%s", author).Blockquote(),
	}
}
