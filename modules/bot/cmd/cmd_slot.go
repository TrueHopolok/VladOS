package cmd

import (
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var fruitWin = []int{}

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
		// scored := getSlotScore(msg.Dice.Value)
		return nil
	},
}

/*
Returns a result of a slot machine in this way:

	slot_values = []string{"bar", "grapes", "lemon", "seven"}

Counts amount of each slot and returns an expected score.
*/
func getSlotScore(value int) int {
	value--
	scores := make([]int, 4)
	for _ = range 3 {
		scores[value%4]++
		value /= 4
	}
	if scores[1] == 2 || scores[2] == 2 {
		return 1
	} else if scores[3] == 2 {
		return 3
	} else if scores[1] == 3 || scores[2] == 3 {
		return 5
	} else if scores[3] == 3 {
		return 10
	} else {
		return 0
	}
}
