package bot

import (
	th "github.com/mymmrac/telego/telegohandler"
)

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
	Handler th.Handler

	// In case the command is multistep (conversation based) this will handle the conversation.
	// See [ConnectConversation] for more details.
	//
	// Value will be nil if conversation handler is not defined.
	Conversation th.Handler
}

// Stores all commands in the map using initialized variables (see [Command] type and its variables).
//
// Few commands are stored and handled seperatly from the list:
//   - [HandleSpelling] is not a command and executed if given command was not spelled correctly (also partially executed during help command, see [HandleHelp]).
//   - [HandleHelp] does not serve any purpose for usage except for guidance, thus stored seperatly (and it has incompability to be stored in global map).
var CommandsList map[string]map[string]Command = map[string]map[string]Command{
	"Gambling": {},
	"Others": {
		"start": CommandStart,
		"ghoul": CommandGhoul,
		"rand":  CommandRand,
	},
}

// Create a group in bot handler that handles all incomming commands.
// See [CommandsList] for all commands details.
func ConnectCommands(bh *th.BotHandler) {
	ch := bh.Group(th.AnyCommand())
	for category := range CommandsList {
		for name, cmd := range CommandsList[category] {
			ch.Handle(cmd.Handler, th.CommandEqual(name))
		}
	}
	ch.Handle(HandleHelp, th.CommandEqual("help"))
	ch.Handle(handleSpelling, th.Any())
}
