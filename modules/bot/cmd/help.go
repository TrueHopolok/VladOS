package cmd

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func handleHelp(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "help")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	_, _, args := tu.ParseCommand(update.Message.Text)
	var msgText []tu.MessageEntityCollection
	if len(args) > 1 {
		msgText = []tu.MessageEntityCollection{tu.Entity("Invalid amount of arguments in the command.\nFor more info type:\n /help <command>\n /help")}
	} else if len(args) == 1 {
		msgText = utilInfoSingle(args[0])
	} else {
		msgText = utilInfoAll()
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
	return err

}

// Returns a message containing info about all of the commands bot has.
func utilInfoAll() []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entity("Here is the whole list of all available commands:\n"))
	for category := range CommandsList {
		msgText = append(msgText, tu.Entityf("\n%s:\n", category).Bold())
		for name, cmd := range CommandsList[category] {
			msgText = append(msgText, tu.Entityf(" /%s - %s\n", name, cmd.InfoBrief))
		}
	}
	msgText = append(msgText, tu.Entity("\nFor more info about particular command use:\n "))
	msgText = append(msgText, tu.Entity("/help <command>").Code())
	return msgText
}

// Returns a message containing a full info about a single command.
func utilInfoSingle(cmdName string) []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	for category := range CommandsList {
		for name, cmd := range CommandsList[category] {
			if name == cmdName {
				msgText = append(msgText, tu.Entityf("%s", cmd.InfoFull))
				return msgText
			}
		}
	}
	msgText = append(msgText, tu.Entityf("'%s' is not a command.\nSee /help for whole list of the commands.\n\nThe most similar commands are:", cmdName))
	for _, potName := range utilClosestSpelling(cmdName) {
		msgText = append(msgText, tu.Entityf("\n /%s", potName))
	}
	return msgText
}
