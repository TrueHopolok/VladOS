package bot

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

var CommandHelp Command = Command{
	Info: `
Has 2 usages:

` + "`/help`" + `
Output a list of all commands in the big message.

` + "`/help <command>`" + `
Output description of the given command.
`,
	Handler:      HandleHelp,
	Conversation: nil,
}

// TODO
func HandleHelp(ctx *th.Context, update telego.Update) error {
	return nil
}
