package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/TrueHopolok/VladOS/modules/db/dbconvo"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type guessConvoStatus struct {
	GuessesLeft  int // lifes left
	PickedNumber int // required to guess
}

const guessesStartingAmount = 6
const guessesMaxValue = 100

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

		_, err = bot.SendMessage(ctx, tu.Messagef(chatID, `
Number was picked, try to guess it.
Type any number between 1 and 100.
I will answer if it is: higher, lower or you guessed it.

If you wish to cancel the guessing type: /cancel
Warning: it will be counted as a loss.

You have 6 guesses, good luck!`))
		if err != nil {
			return fmt.Errorf("send msg: %w", err)
		}

		status := guessConvoStatus{GuessesLeft: guessesStartingAmount, PickedNumber: rand.New(rand.NewSource(time.Now().UnixNano())).Intn(guessesMaxValue) + 1}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(status); err != nil {
			return fmt.Errorf("gob encoder: %w", err)
		}
		return dbconvo.Busy(update.Message.From.ID, "guess", buf.Bytes())
	},
	conversation: func(ctx *telegohandler.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "guess")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		userID := update.Message.From.ID

		cs := ctx.Value(ctxValueConvoStatus{}).(dbconvo.Status)
		getbuf := bytes.NewBuffer(cs.Data)
		dec := gob.NewDecoder(getbuf)
		var status guessConvoStatus
		if err := dec.Decode(&status); err != nil {
			return fmt.Errorf("gob decoder: %w", err)
		}
		status.GuessesLeft--

		playerGuess, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Given argument is invalid, please enter the valid number between 1 and %d (included).\nTo cancel the command input and execution type:\n /cancel", guessesMaxValue))
			return err
		}
		if playerGuess == status.PickedNumber {
			msgText, err := utilOutputDice("guess", userID, true)
			if err != nil {
				return err
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
			if err != nil {
				return err
			}
			return dbconvo.Free(userID)
		}

		if status.GuessesLeft == 0 {
			msgText, err := utilOutputDice("guess", userID, false)
			if err != nil {
				return err
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
			if err != nil {
				return err
			}
			return dbconvo.Free(userID)
		}

		var answer string
		if playerGuess > status.PickedNumber {
			answer = "LOWER"
		} else {
			answer = "GREATER"
		}
		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, tu.Entityf("You guessed wrong.\n%d guesses left.\n\nPicked number is ", status.GuessesLeft), tu.Entityf("%s the yours.", answer).Bold()))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(status); err != nil {
			return fmt.Errorf("gob encoder: %w", err)
		}
		return dbconvo.Busy(update.Message.From.ID, "guess", buf.Bytes())
	},
	cancelation: func(ctx *telegohandler.Context, update telego.Update) error {
		return dbstats.Update("guess", update.Message.From.ID, 0)
	},
}
