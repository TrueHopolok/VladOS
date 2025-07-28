package bot

import (
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
		bot, chatID, _, valid, err := CmdStart(ctx, update, "ghoul", 0)
		if !valid {
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
