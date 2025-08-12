package cmd

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: make a blackjack game
var CommandBjack Command = Command{
	InfoBrief: "game of a blackjack",
	InfoFull: `
 /bjack
Play a blackjack against a classical dealer. 
No need to bet any money, since you will bet your score streak like in dice and slots.

Gameplay:
Sorry, but this command is currently in development.

On losing score is reset.

Has a leaderboard to count largest score streak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "blackjack", 0)
		if !valid {
			return err
		}

		_, err = bot.SendMessage(ctx, tu.Message(chatID, "Sorry, but currently this command is in development. Comeback later."))
		return err
	},
	conversation: nil,
	cancelation:  nil,
}
