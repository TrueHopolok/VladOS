package cmd

import (
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
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
		case "slot", "dice", "bjack":
			ranking, err := dbstats.GetFull(args[0], userID)
			if err != nil {
				return err
			}
			if len(ranking) == 0 {
				_, err = bot.SendMessage(ctx, tu.Message(chatID, "Currently no players have played this game. Become the first one."))
				return err
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, outputStats(userID, ranking)...))
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
}

func outputStats(userID int64, ranking []dbstats.FullStats) []tu.MessageEntityCollection {
	var yourStats dbstats.UserStats
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
		msgText = append(msgText, tu.Entity("\n\nYour stats:").Bold(), tu.Entityf("\nGames in total: %d\nCurrent score streak: %d\nBest score streak: %d", yourStats.GamesTotal, yourStats.ScoreCurrent, yourStats.ScoreBest))
	} else {
		msgText = append(msgText, tu.Entity("\n\nPlay the game to have any stats..."))
	}
	return msgText
}
