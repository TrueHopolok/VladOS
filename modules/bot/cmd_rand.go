package bot

import (
	"fmt"

	// dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const RandMaxValue int = 1_000_000_000

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
			return err
		}
		// return dbconvo.Busy(update.Message.From.ID, "rand", object with stage and left border info -> cs.Data)
		return nil
	},
	Conversation: func(ctx *th.Context, update telego.Update) error {
		// cs := ctx.Value("ConversationStatus").(dbconvo.Status)
		// cs.Data -> object with stage and left border info
		return nil
	},
}
