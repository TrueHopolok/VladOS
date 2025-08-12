package cmd

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: make a blackjack game
var CommandGuess Command = Command{
	InfoBrief: "game of a number guesser",
	InfoFull: `
 /guess
Play a number guesser game against the bot in a role of the guesser.
He randomly picks a number between 1 and 100.
Your goal is to guess it in 6 or less guesses. 

On correct guess your score streak is increased.
On losing score is reset.

Has a leaderboard to count largest score streak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "guess", 0)
		if !valid {
			return err
		}

		_, err = bot.SendMessage(ctx, tu.Message(chatID, "Sorry, but currently this command is in development. Comeback later."))
		return err
	},
	conversation: nil,
	cancelation:  nil,
}
