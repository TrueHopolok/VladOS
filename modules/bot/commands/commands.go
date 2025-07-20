// Connector package for bot handlers to make it easier to use.
//
// Main funciton to use is [Connect]. See it or [CommandsList] for more info.
package commands

import (
	th "github.com/mymmrac/telego/telegohandler"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Contains the command's handler and the info text that
// is displayed on help command.
//
// Name of the command is stored in key of the [CommandsList] map.
type Command struct {
	Info    string
	Handler th.Handler
}

// Stores all commands in the map using Get functions.
// Get functions are wrappers for handler functions to create a Command with given info (see [GetHelp] and [HandleHelp]).
//
// Few commands are stored and handled seperatly from the list:
//   - [HandleSpelling] is not a command and executed if given command was not spelled correctly (also partially executed during help command see [HandleHelp]).
//   - [HandleStart] should be used only once on initialization, thus is executed seperatly.
var CommandsList map[string]Command = map[string]Command{
	"help": GetHelp(),
}

// Create a group in bot handler that handles all incomming commands.
// See [CommandsList] for all commands details.
func Connect(bh *th.BotHandler) {
	ch := bh.Group(th.AnyCommand())
	ch.Handle(HandleStart, th.CommandEqual("start"))
	for name, cmd := range CommandsList {
		ch.Handle(cmd.Handler, th.CommandEqual(name))
	}
	ch.Handle(HandleSpelling, th.Any())
}
