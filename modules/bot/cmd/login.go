package cmd

import (
	"github.com/TrueHopolok/VladOS/modules/db/dblogin"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandLogin Command = Command{
	InfoBrief: "generate auth code",
	InfoFull: `
 /login
Used to generate authefication code.
It is required to login on the web page.

Code is valid for 5 minutes, afterwards it will be deleted.
`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "login", 0)
		if !valid {
			return err
		}

		authcode, err := dblogin.Add(update.Message.From.ID, update.Message.From.FirstName, update.Message.From.Username)
		if err != nil {
			return err
		}

		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID,
			tu.Entity("Generating authefication code successeded.\n\n"),
			tu.Entity("DO NOT SHARE IT\nIT EXPIRES AFTER 5 MINUTES").Bold(),
			tu.Entity("\n\nYour code is:\n"),
			tu.Entity(authcode).Code(),
		))
		return err
	},
}
