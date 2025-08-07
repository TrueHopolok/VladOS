package cmd

import (
	"fmt"

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
		// msg.Dice.Value == 6
		return nil
	},
}
