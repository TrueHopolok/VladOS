package cmd

import (
	"math/rand"
	"time"

	"github.com/TrueHopolok/VladOS/modules/db/dbm8b"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var CommandM8B Command = Command{
	InfoBrief: "magic 8 ball",
	InfoFull: `
 /mb
Tells you the truth to a yes/no question in mind.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "mb", 0)
		if !valid {
			return err
		}

		ans_int := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2)
		ans_str, err := dbm8b.Get(ans_int == 1)
		if err != nil {
			return err
		}

		_, err = bot.SendMessage(ctx, tu.Message(chatID, ans_str))
		return err
	},
	conversation: nil,
}
