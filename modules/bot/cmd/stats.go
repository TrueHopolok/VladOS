package cmd

import (
	"github.com/TrueHopolok/VladOS/modules/db/dbslot"
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

		userID := update.Message.From.ID
		switch args[0] {
		case "slot":
			ranking, err := dbslot.GetFull(userID)
			if err != nil {
				return err
			}
			if len(ranking) == 0 {
				_, err = bot.SendMessage(ctx, tu.Message(chatID, "Currently no players have played this game. Become the first one."))
				return err
			}
			var yourStats dbslot.UserStats
			foundYou := false
			msgText := []tu.MessageEntityCollection{
				tu.Entityf("Players in total:\n%d\n\nLeaderboard:", ranking[0].PlayersTotal).Bold(),
			}
			for _, stats := range ranking {
				msgText = append(msgText, tu.Entityf("\n%d. place: %d", stats.Placement, stats.Personal.ScoreBest))
				if stats.UserId == userID {
					foundYou = true
					yourStats = stats.Personal
					msgText = append(msgText, tu.Entity("  (you)").Bold())
				}
			}
			if foundYou {
				msgText = append(msgText, tu.Entity("\n\nYour stats:").Bold(), tu.Entityf("\nTotal spins: %d\nCurrent score streak: %d\nBest score streak: %d", yourStats.SpinsTotal, yourStats.ScoreCurrent, yourStats.ScoreBest))
			} else {
				msgText = append(msgText, tu.Entity("\n\nPlay the game to have any stats..."))
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
			return err
		case "dice":
			_, err = bot.SendMessage(ctx, tu.Message(chatID, "Sorry, but currently this command is in development. Comeback later."))
			return err
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
