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
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const randMaxValue int = 1_000_000_000

type randConvoStatus struct {
	// 0 - start, waiting left one
	// 1 - left is given, waiting right one
	Stage int

	Left int
}

var CommandRand Command = Command{
	InfoFull: fmt.Sprintf(`
 /rand
Generates a random number between given numbers in range from 0 till %d included.
Command will give you a prompt asking what min and max values of the random you want.

 /rand <max_num>
Generates a random number between 0 and given number which value is from 0 till %d included.
Command will immediatly send a response. Expects that max_num is in the allowed range.

 /rand <min_num> <max_num>
Generates a random number between 2 given numbers which values is from 0 till %d included.
Command will immediatly send a response. Expects that min_num <= max_num and they are in the allowed range.
`, randMaxValue, randMaxValue, randMaxValue),
	InfoBrief: "generates random number",
	handler: func(ctx *th.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "rand")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		_, _, cmdArgs := tu.ParseCommand(update.Message.Text)
		switch len(cmdArgs) {
		case 0:
			_, err := bot.SendMessage(ctx, tu.Messagef(chatID, "Type what minimum value is allowed.\nAllowed values are between 0 and %d (included).", randMaxValue))
			if err != nil {
				return fmt.Errorf("send msg: %w", err)
			}
			status := randConvoStatus{
				Stage: 0,
				Left:  0,
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(status); err != nil {
				return fmt.Errorf("gob encoder: %w", err)
			}
			return dbconvo.Busy(update.Message.From.ID, "rand", buf.Bytes())
		case 1:
			right, invalid, err := inputRand(ctx, chatID, cmdArgs[0], 0, false)
			if invalid {
				return err
			}

			return outputRand(ctx, chatID, 0, right)
		case 2:
			left, invalid, err := inputRand(ctx, chatID, cmdArgs[0], 0, false)
			if invalid {
				return err
			}

			right, invalid, err := inputRand(ctx, chatID, cmdArgs[1], left, false)
			if invalid {
				return err
			}

			return outputRand(ctx, chatID, left, right)
		default:
			_, err := bot.SendMessage(ctx, tu.Message(chatID, "Too many arguments are given for the command.\nFor more info type:\n /help rand\n /help"))
			return err
		}
	},
	conversation: func(ctx *th.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "rand")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		cs := ctx.Value(ctxValueConvoStatus{}).(dbconvo.Status)
		getbuf := bytes.NewBuffer(cs.Data)
		dec := gob.NewDecoder(getbuf)
		var status randConvoStatus
		if err := dec.Decode(&status); err != nil {
			return fmt.Errorf("gob decoder: %w", err)
		}
		switch status.Stage {
		case 0: // start, waiting left one
			left, invalid, err := inputRand(ctx, chatID, update.Message.Text, 0, true)
			if invalid {
				return err
			}

			_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Type what maximum value is allowed.\nAllowed values are between %d and %d (included).", left, randMaxValue))
			if err != nil {
				return fmt.Errorf("send msg: %w", err)
			}

			status.Left = left
			status.Stage++
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(status); err != nil {
				return fmt.Errorf("gob encoder: %w", err)
			}
			return dbconvo.Busy(update.Message.From.ID, "rand", buf.Bytes())
		case 1: // left is given, waiting right one
			right, invalid, err := inputRand(ctx, chatID, update.Message.Text, status.Left, true)
			if invalid {
				return err
			}

			if err := outputRand(ctx, chatID, status.Left, right); err != nil {
				return err
			}

			return dbconvo.Free(update.Message.From.ID)
		}
		return nil
	},
	cancelation: func(ctx *th.Context, update telego.Update) error {
		return nil
	},
}

func inputRand(ctx *th.Context, chatID telego.ChatID, inputText string, left int, withCancel bool) (inputNum int, invalid bool, msgErr error) {
	bot := ctx.Bot()
	num, err := strconv.Atoi(inputText)
	if err != nil || num < left || num > randMaxValue {
		if withCancel {
			_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Given argument is invalid, please enter the valid number between %d and %d (included).\nTo cancel the command input and execution type:\n /cancel", left, randMaxValue))
		} else {
			_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Given argument is invalid, please enter the valid number between %d and %d (included).\nFor more info type:\n /help rand\n /help", left, randMaxValue))
		}
		return 0, true, err
	}
	return num, false, nil
}

func outputRand(ctx *th.Context, chatID telego.ChatID, left, right int) error {
	bot := ctx.Bot()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err := bot.SendMessage(ctx, tu.Messagef(chatID, "Generated number between %d and %d is:", left, right))
	if err != nil {
		return err
	}
	_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, tu.Entityf("%d", r.Intn(right-left+1)+left)))
	return err
}
