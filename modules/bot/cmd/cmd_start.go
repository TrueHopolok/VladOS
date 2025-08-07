package cmd

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandStart Command = Command{
	InfoBrief: "output info about the bot and its links",
	InfoFull: `
 /start
Output information about bot's life purpose and who it is.
Provide useful links to bot's website. 
`,
	handler: func(ctx *th.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "start", 0)
		if !valid {
			return err
		}
		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID,
			tu.Entity("Hello, "), tu.Entity("user.\n").Bold(),
			tu.Entity("I am bot "), tu.Entity("VladOS.\nVlad Operation System.\n").Bold(),
			tu.Entity(`
I am a project that combines:
- Gambling telegram bot;
- Reincarnaction of the 'AllEgg' bot from Discord;
- Webpage to control and view bot activities.

Bot will react to any message and will try to find a pun for suffix of the message.

Gambling and other functional via commands.
Type /help for more info about them.
`)))
		return err
	},
	conversation: nil,
}
