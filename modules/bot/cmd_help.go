package bot

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: add seperation between empty argument command and non empty argument
func HandleHelp(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "help")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	_, _, args := tu.ParseCommand(update.Message.Text)
	var msgText []tu.MessageEntityCollection
	if len(args) > 1 {
		msgText = CmdInvalidArgsAmount()
	} else if len(args) == 1 {
		msgText = CmdInfoOne(args[0])
	} else {
		msgText = CmdInfoAll()
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
	return err

}

// Create a message to send about invalid amount of arguments in the command.
func CmdInvalidArgsAmount() []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entity("Invalid amount of arguments in the command.\nTry /help or /help <command>."))
	return msgText
}

// Returns a message containing info about all of the commands bot has.
func CmdInfoAll() []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entity("Here is the whole list of all available commands:\n\n"))
	for name, cmd := range CommandsList {
		msgText = append(msgText, tu.Entityf(" /%s - %s\n", name, cmd.InfoBrief))
	}
	msgText = append(msgText, tu.Entity("\nFor more info about particular command use:\n /help "))
	msgText = append(msgText, tu.Entity("<command>").Code())
	return msgText
}

// Returns a message containing a full info about a single command.
func CmdInfoOne(cmdName string) []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entityf("%s", CommandsList[cmdName].InfoFull))
	return msgText
}
