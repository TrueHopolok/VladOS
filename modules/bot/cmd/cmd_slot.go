package cmd

import (
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandSlot Command = Command{
	InfoBrief: "spin slot machine",
	InfoFull: `
 /slot
Spin a slot machine. Win is recieving a Jackpot (1:64 chance).

Has a leaderboard to count largest winstreak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "dice", 0)
		if !valid {
			return err
		}

		msg, err := bot.SendDice(ctx, tu.Dice(chatID, telego.EmojiSlotMachine))
		if err != nil {
			return err
		}

		if msg.Dice == nil {
			return fmt.Errorf("msg is not a dice result: %v", msg)
		}
		// msg.Dice.Value == 64
		return nil
	},
}
