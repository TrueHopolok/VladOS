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
	Info string

	// Command handler that executes on command call.
	Handler th.Handler

	// In case the command is multistep (conversation based) this will handle the conversation.
	// See [ConnectConversation] for more details.
	//
	// Value will be nil in case conversation is not intended.
	Conversation *th.Handler
}

// Stores all commands in the map using initialized variables (see [Command] type and its variables).
//
// Few commands are stored and handled seperatly from the list:
//   - [HandleSpelling] is not a command and executed if given command was not spelled correctly (also partially executed during help command, see [HandleHelp]).
//   - [HandleHelp] does not serve any purpose for usage except for guidance, thus stored seperatly.
//   - [HandleCancel] is used in conversation only, thus is not a independed command.
//   - [HandleStart] should be used once thus no need to include in the whole command list.
var CommandsList map[string]Command = map[string]Command{}

// Create a group in bot handler that handles all incomming commands.
// See [CommandsList] for all commands details.
func ConnectCommands(bh *th.BotHandler) {
	ch := bh.Group(th.AnyCommand())
	for name, cmd := range CommandsList {
		ch.Handle(cmd.Handler, th.CommandEqual(name))
	}
	ch.Handle(HandleStart, th.CommandEqual("start"))
	ch.Handle(HandleHelp, th.CommandEqual("help"))
	ch.Handle(HandleSpelling, th.Any())
}
