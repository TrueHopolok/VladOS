package bot

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: add seperation between empty argument command and non empty argument
func HandleHelp(ctx *th.Context, update telego.Update) error {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entity("Here is the whole list of all available commands:\n\n"))
	for name := range CommandsList {
		msgText = append(msgText, tu.Entityf("/%s\n", name))
	}
	msgText = append(msgText, tu.Entity("\nFor more info about particular command use:\n /help "))
	msgText = append(msgText, tu.Entity("<command>").Bold())
	bot := ctx.Bot()
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(update.Message.Chat.ChatID(), msgText...))
	return err
}
