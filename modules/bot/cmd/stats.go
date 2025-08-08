package cmd

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: make whole stats and update on for each new machine (potentialy make it dynamic)
var CommandStats Command = Command{
	InfoBrief: "output stats for game",
	InfoFull: `
 /stats <game_name>
Output all stats for ceratin game.
With a your placement and top placement in the leaderboard.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, args, valid, err := utilStart(ctx, update, "stats", 1)
		if !valid {
			return err
		}

		_, err = bot.SendMessage(ctx, tu.Message(chatID, "Sorry, but currently this command is in development. Comeback later."))
		return err

		switch args[0] {
		case "slot":
			return nil
		case "dice":
			return nil
		case "bjack":
			_, err = bot.SendMessage(ctx, tu.Message(chatID, "Sorry, but currently this command is in development. Comeback later."))
			return err
		default:
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID,
				tu.Entity("Provided categories does not exists.\nTry any of those:\n"),
				tu.Entity(" /stats slot\n").Code(),
				tu.Entity(" /stats dice\n").Code(),
				tu.Entity(" /stats bjack\n").Code()))
			return err
		}
	},
	conversation: nil,
}
