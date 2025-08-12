package cmd

import (
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandDice Command = Command{
	InfoBrief: "throw a dice",
	InfoFull: `
 /dice
Throw a dice. Win is hitting anything except one (5:6 chance).
The value of the dice is counted towards score.
On losing score is reset.

Has a leaderboard to count largest score streak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "dice", 0)
		if !valid {
			return err
		}

		msg, err := bot.SendDice(ctx, tu.Dice(chatID, telego.EmojiDice))
		if err != nil {
			return err
		}

		if msg.Dice == nil {
			return fmt.Errorf("msg is not a dice result: %v", msg)
		}
		diceScore := msg.Dice.Value
		if diceScore == 1 {
			diceScore = 0
		}
		err = dbstats.Update("dice", update.Message.From.ID, msg.Dice.Value)
		if err != nil {
			return err
		}
		msgText, err := utilOutputDice("dice", update.Message.From.ID, msg.Dice.Value > 0)
		if err != nil {
			return err
		}
		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
		return err
	},
}

func utilOutputDice(gameName string, userID int64, hasWon bool) ([]tu.MessageEntityCollection, error) {
	stats, err := dbstats.Get(gameName, userID)
	if err != nil {
		return nil, err
	}

	var msgText []tu.MessageEntityCollection
	if hasWon {
		msgText = append(msgText, tu.Entity("You won!\n").Bold())
	} else {
		msgText = append(msgText, tu.Entity("You lost!\n").Bold())
	}

	msgText = append(msgText, tu.Entityf("\nCurrent score: %d\nBest score: %d\nMore stats:", stats.ScoreCurrent, stats.ScoreBest), tu.Entityf(" /stats %s", gameName).BotCommand(), tu.Entityf("\nPlay again: /%s", gameName))
	return msgText, nil
}
