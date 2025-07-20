package commands

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

//	Command{
//		Info: 	"Information about the help command...",
//		Handler: HandleHelp,
//	}
func GetHelp() Command {
	return Command{
		Info: `
Has 2 usages:

` + "`/help`" + `
Output a list of all commands in the big message.

` + "`/help <command>`" + `
Output description of the given command.
`,
		Handler: HandleHelp,
	}
}

// TODO
func HandleHelp(ctx *th.Context, update telego.Update) error {
	return nil
}
