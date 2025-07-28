package bot

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleHelp(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "help")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	_, _, args := tu.ParseCommand(update.Message.Text)
	var msgText []tu.MessageEntityCollection
	if len(args) > 1 {
		msgText = []tu.MessageEntityCollection{tu.Entity("Invalid amount of arguments in the command.\nFor more info type:\n /help <command>\n /help")}
	} else if len(args) == 1 {
		msgText = CmdInfoOne(args[0])
	} else {
		msgText = CmdInfoAll()
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
	return err

}

// Outputs log with info [cmdName] in the [log/slog].
// Checks if received [len(args)] is equal to given [argsAmount].
// Sends a message if it is false.
func CmdStart(ctx *th.Context, update telego.Update, cmdName string, argsAmount int) (bot *telego.Bot, chatID telego.ChatID, cmdArgs []string, validArgs bool, invalidMSG error) {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", cmdName)
	bot = ctx.Bot()
	chatID = update.Message.Chat.ChatID()
	_, _, cmdArgs = tu.ParseCommand(update.Message.Text)
	validArgs = len(cmdArgs) == argsAmount
	if !validArgs {
		_, invalidMSG = bot.SendMessage(ctx, tu.Messagef(chatID, "Invalid amount of arguments in the command.\nFor more info type:\n /help %s\n /help", cmdName))
	}
	return bot, chatID, cmdArgs, validArgs, invalidMSG
}

// Returns a message containing info about all of the commands bot has.
func CmdInfoAll() []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entity("Here is the whole list of all available commands:\n"))
	for category := range CommandsList {
		msgText = append(msgText, tu.Entityf("\n%s:\n", category).Bold())
		for name, cmd := range CommandsList[category] {
			msgText = append(msgText, tu.Entityf(" /%s - %s\n", name, cmd.InfoBrief))
		}
	}
	msgText = append(msgText, tu.Entity("\nFor more info about particular command use:\n /help "))
	msgText = append(msgText, tu.Entity("<command>").Code())
	return msgText
}

// Returns a message containing a full info about a single command.
func CmdInfoOne(cmdName string) []tu.MessageEntityCollection {
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
	for _, potName := range CmdFindClosest(cmdName) {
		msgText = append(msgText, tu.Entityf("\n /%s", potName))
	}
	return msgText
}
