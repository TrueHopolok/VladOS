package cmd

import (
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db/dbslot"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandSlot Command = Command{
	InfoBrief: "spin slot machine",
	InfoFull: `
 /slot
Spin a slot machine.
Win is counted if two or three of 1 kind with excpetion of bar.
Score board (with its chances):
 - Two of a kind   =  +1 (27:64);
 - Double sevens   =  +3 (9:64);
 - Three of a kind   =  +9 (3:64);
 - Triple sevens / JACKPOT   = +27 (1:64).

Overall chance to win and continue the streak is (40:64) or (5:8).
On losing score is reset.

Has a leaderboard to count largest score streak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "slot", 0)
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
		scored := getSlotScore(msg.Dice.Value)
		err = dbslot.Update(update.Message.From.ID, scored)
		if err != nil {
			return err
		}
		msgText, err := outputSlot(update.Message.From.ID, scored > 0)
		if err != nil {
			return err
		}
		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
		return err
	},
}

func outputSlot(user_id int64, has_won bool) ([]tu.MessageEntityCollection, error) {
	stats, err := dbslot.Get(user_id)
	if err != nil {
		return nil, err
	}

	var msgText []tu.MessageEntityCollection
	if has_won {
		msgText = append(msgText, tu.Entity("You won!\n").Bold())
	} else {
		msgText = append(msgText, tu.Entity("You lost!\n").Bold())
	}

	msgText = append(msgText, tu.Entityf("\nCurrent score: %d\nBest score: %d\nMore stats:", stats.ScoreCurrent, stats.ScoreBest), tu.Entity(" /stats slot").BotCommand(), tu.Entity("\nPlay again: /slot"))
	return msgText, nil
}

/*
Returns a result of a slot machine in this way:

	slot_values = []string{"bar", "grapes", "lemon", "seven"}

Counts amount of each slot and returns an expected score.
*/
func getSlotScore(value int) int {
	value--
	scores := make([]int, 4)
	for range 3 {
		scores[value%4]++
		value /= 4
	}

	switch 2 {
	case scores[0], scores[1], scores[2]:
		return 1
	case scores[3]:
		return 3
	}

	switch 3 {
	case scores[0], scores[1], scores[2]:
		return 9
	case scores[3]:
		return 27
	}

	return 0
}
