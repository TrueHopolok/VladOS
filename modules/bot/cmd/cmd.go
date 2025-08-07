// All bot commands are handled here.
// Misspelling is handled by github.com/TrueHopolok/spellchecker/spellchecker.
//
// Command handling is processed in this order:
//
//	get conversation status
//	if ( status is active ):
//		if msg is cmd{cancel}: reset conversation status
//		else: give control to cmd converstation handler
//	else if ( msg is cmd ):
//		if ( msg in list{cmd} ): give control to cmd handler
//		else ( if msg is cmd{help} ): output info about commands and bot
//	 	else: handle that as cmd misspell
//	else: skip this package and fallthrough to others handlers
package cmd

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Contains the command's handler and the info text that
// is displayed on help command.
//
// Name of the command is stored in key of the [CommandsList] map.
type Command struct {
	// Description and a usage of the command.
	InfoFull string

	// Brief description about the command.
	InfoBrief string

	// Command handler that executes on command call.
	handler th.Handler

	// In case the command is multistep (conversation based) this will handle the conversation.
	// See [ConnectConversation] for more details.
	//
	// Value will be nil if conversation handler is not defined.
	conversation th.Handler
}

// Stores all commands in the map using initialized variables (see [Command] type and its variables).
//
// Few commands are stored and handled seperatly from the list:
//   - help command since it require accessing this map, which is impossible wihtout cyclic import
//   - spell checking since it is not a command and happens in case of misspelled command
var CommandsList map[string]map[string]Command = map[string]map[string]Command{
	"Gambling": {},
	"Others": {
		"start": CommandStart,
		"ghoul": CommandGhoul,
		"rand":  CommandRand,
		"tip":   CommandTip,
	},
}

// Connects converstion handlers.
// Afterwards create a group in bot handler that handles all incomming commands.
// See [CommandsList] for all commands details.
func ConnectCommands(bh *th.BotHandler) error {
	if err := connectConversation(bh); err != nil {
		return err
	}
	ch := bh.Group(th.AnyCommand())
	for category := range CommandsList {
		for name, cmd := range CommandsList[category] {
			ch.Handle(cmd.handler, th.CommandEqual(name))
		}
	}
	ch.Handle(handleHelp, th.CommandEqual("help"))
	ch.Handle(handleSpelling, th.AnyMessage())
	return nil
}

// Outputs log with info [cmdName] in the [log/slog].
// Checks if received [len(args)] is equal to given [argsAmount].
// Sends a message if it is false.
func utilStart(ctx *th.Context, update telego.Update, cmdName string, argsAmount int) (bot *telego.Bot, chatID telego.ChatID, cmdArgs []string, validArgs bool, invalidMSG error) {
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
