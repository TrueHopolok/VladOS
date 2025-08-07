package cmd

import (
	"fmt"

	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type ctxValueConvoStatus struct{}

func handleConvoStatus(ctx *th.Context, update telego.Update) error {
	cs, err := dbconvo.Get(update.Message.From.ID)
	if err != nil {
		return err
	}
	ctx = ctx.WithValue(ctxValueConvoStatus{}, cs)
	return ctx.Next(update)
}

func handleConvoCancel(ctx *th.Context, update telego.Update) error {
	bot, chatID, _, valid, err := utilStart(ctx, update, "cancel", 0)
	if !valid {
		return err
	}
	_, err = bot.SendMessage(ctx, tu.Message(chatID, "Conversation/command was canceled.\nBot returned to normal state."))
	if err != nil {
		return err
	}
	return dbconvo.Free(update.Message.From.ID)
}

func handleConversation(ctx *th.Context, update telego.Update) error {
	cs := ctx.Value(ctxValueConvoStatus{}).(dbconvo.Status)
	if cs.Available {
		return ctx.Next(update)
	}

	for category := range CommandsList {
		cmd, exists := CommandsList[category][cs.CommandName]
		if !exists {
			continue
		} else if cmd.conversation == nil {
			return fmt.Errorf("%s command conversation handler is nil, but got in conversation status", cs.CommandName)
		}
		return (cmd.conversation)(ctx, update)
	}
	return fmt.Errorf("%s command does not exists in the commands list, but got in conversation status", cs.CommandName)
}

func connectConversation(bh *th.BotHandler) error {
	if err := dbconvo.Clear(); err != nil {
		return err
	}
	bh.Handle(handleConvoStatus, th.AnyMessage())
	bh.Handle(handleConvoCancel, th.CommandEqual("cancel"))
	bh.Handle(handleConversation, th.AnyMessage())
	return nil
}
