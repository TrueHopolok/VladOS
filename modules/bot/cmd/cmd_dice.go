package cmd

import (
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db/dice"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandDice Command = Command{
	InfoBrief: "throw a dice",
	InfoFull: `
 /dice
Throw a dice. Win is recieving hitting a six (1:6 chance).

Has a leaderboard to count largest winstreak.`,
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
		err = dice.Update(update.Message.From.ID, update.Message.From.FirstName, msg.Dice.Value)
		if err != nil {
			return err
		}
		msgText, err := outputDice(update.Message.From.ID, msg.Dice.Value == 6)
		if err != nil {
			return err
		}
		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
		return err
	},
}

func outputDice(user_id int64, has_won bool) ([]tu.MessageEntityCollection, error) {
	stats, err := dice.Get(user_id)
	if err != nil {
		return nil, err
	}

	var msgText []tu.MessageEntityCollection
	if has_won {
		msgText = append(msgText, tu.Entity("You won!\n").Bold())
	} else {
		msgText = append(msgText, tu.Entity("You lost!\n").Bold())
	}

	msgText = append(msgText, tu.Entityf("\nCurrent winstreak: %d\nBest winstreak: %d\nMore stats:", stats.StreakCurrent, stats.StreakBest), tu.Entity(" /stats dice").BotCommand(), tu.Entity("\nRepeat: /dice"))
	return msgText, nil
}
