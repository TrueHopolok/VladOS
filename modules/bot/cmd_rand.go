package bot

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const RandMaxValue int = 1_000_000_000

type CmdRandStatus struct {
	// 0 - start, waiting left one
	// 1 - left is given, waiting right one
	Stage int

	Left int
}

var CommandRand Command = Command{
	InfoFull: fmt.Sprintf(`
 /rand
Generates a random number between 0 and %d including.
Command will give you a prompt asking what min and max values of the random you want.
`, RandMaxValue),
	InfoBrief: "generates random number",
	Handler: func(ctx *th.Context, update telego.Update) error {
		bot, chatID, _, valid, err := CmdStart(ctx, update, "rand", 0)
		if !valid {
			return err
		}
		_, err = bot.SendMessage(ctx, tu.Message(chatID, "Type what minimum value (including) is allowed."))
		if err != nil {
			return fmt.Errorf("send msg: %w", err)
		}
		status := CmdRandStatus{
			Stage: 0,
			Left:  0,
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(status); err != nil {
			return fmt.Errorf("gob encoder: %w", err)
		}
		return dbconvo.Busy(update.Message.From.ID, "rand", buf.Bytes())
	},
	Conversation: func(ctx *th.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "rand")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		cs := ctx.Value("ConvoStatus").(dbconvo.Status)
		getbuf := bytes.NewBuffer(cs.Data)
		dec := gob.NewDecoder(getbuf)
		var status CmdRandStatus
		if err := dec.Decode(&status); err != nil {
			return fmt.Errorf("gob decoder: %w", err)
		}
		switch status.Stage {
		case 0: // start, waiting left one
			value, err := strconv.Atoi(update.Message.Text)
			if err != nil || value < 0 || value > RandMaxValue {
				_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Inputed text is invalid number, please enter the valid number between 0 and %d (included).\nTo cancel execution of the command type /cancel.", RandMaxValue))
				if err != nil {
					return fmt.Errorf("send msg: %w", err)
				}
				return nil
			}

			_, err = bot.SendMessage(ctx, tu.Message(chatID, "Type what maximum value (including) is allowed."))
			if err != nil {
				return fmt.Errorf("send msg: %w", err)
			}

			status.Left = value
			status.Stage++
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(status); err != nil {
				return fmt.Errorf("gob encoder: %w", err)
			}
			return dbconvo.Busy(update.Message.From.ID, "rand", buf.Bytes())
		case 1: // left is given, waiting right one
			value, err := strconv.Atoi(update.Message.Text)
			if err != nil || value < status.Left || value > RandMaxValue {
				_, err = bot.SendMessage(ctx, tu.Messagef(chatID, "Inputed text is invalid number, please enter the valid number between %d and %d (included).\nTo cancel execution of the command type /cancel.", status.Left, RandMaxValue))
				if err != nil {
					return fmt.Errorf("send msg: %w", err)
				}
				return nil
			}

			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, tu.Entityf("Generated number between %d and %d is:\n\n", status.Left, value), tu.Entityf("%d", r.Intn(value+1)+status.Left).Blockquote()))
			if err != nil {
				return fmt.Errorf("send msg: %w", err)
			}

			return dbconvo.Free(update.Message.From.ID)
		}
		return nil
	},
}
