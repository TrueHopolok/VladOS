package bot

import (
	"fmt"

	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func HandleConversation(ctx *th.Context, update telego.Update) error {
	cs, err := dbconvo.Get(update.Message.From.ID)
	if err != nil {
		return err
	} else if cs.Available {
		return ctx.Next(update)
	}

	cmd, isGambling := CommandsList.Gambling[cs.CommandName]
	cmd, isOthers := CommandsList.Others[cs.CommandName]
	if !isGambling && !isOthers {
		return fmt.Errorf("%s command does not exists in the commands list, but got in conversation status", cs.CommandName)
	} else if cmd.Conversation == nil {
		return fmt.Errorf("%s command conversation handler is nil, but got in conversation status", cs.CommandName)
	}
	return (*(cmd.Conversation))(ctx, update)
}

func ConnectConversation(bh *th.BotHandler) {
	bh.Handle(HandleConversation, th.AnyMessage())
}
