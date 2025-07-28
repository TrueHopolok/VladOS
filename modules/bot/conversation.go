package bot

import (
	"fmt"

	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleConvoStatus(ctx *th.Context, update telego.Update) error {
	cs, err := dbconvo.Get(update.Message.From.ID)
	if err != nil {
		return err
	}
	ctx = ctx.WithValue("ConvoStatus", cs)
	return ctx.Next(update)
}

func HandleConvoCancel(ctx *th.Context, update telego.Update) error {
	bot, chatID, _, valid, err := CmdStart(ctx, update, "cancel", 0)
	if !valid {
		return err
	}
	_, err = bot.SendMessage(ctx, tu.Message(chatID, "Conversation/command was canceled.\nBot returned to normal state."))
	if err != nil {
		return err
	}
	return dbconvo.Free(update.Message.From.ID)
}

func HandleConversation(ctx *th.Context, update telego.Update) error {
	cs := ctx.Value("ConvoStatus").(dbconvo.Status)
	if cs.Available {
		return ctx.Next(update)
	}

	for category := range CommandsList {
		cmd, exists := CommandsList[category][cs.CommandName]
		if !exists {
			continue
		} else if cmd.Conversation == nil {
			return fmt.Errorf("%s command conversation handler is nil, but got in conversation status", cs.CommandName)
		}
		return (cmd.Conversation)(ctx, update)
	}
	return fmt.Errorf("%s command does not exists in the commands list, but got in conversation status", cs.CommandName)
}

func ConnectConversation(bh *th.BotHandler) {
	bh.Handle(HandleConvoStatus, th.AnyMessage())
	bh.Handle(HandleConvoCancel, th.CommandEqual("cancel"))
	bh.Handle(HandleConversation, th.AnyMessage())
}
