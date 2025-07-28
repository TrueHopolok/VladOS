package bot

import (
	"fmt"
	"log/slog"

	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStatus(ctx *th.Context, update telego.Update) error {
	cs, err := dbconvo.Get(update.Message.From.ID)
	if err != nil {
		return err
	}
	ctx = ctx.WithValue("ConversationStatus", cs)
	return ctx.Next(update)
}

func HandleCancel(ctx *th.Context, update telego.Update) error {
	cs := (ctx.Value("ConversationStatus")).(dbconvo.Status)
	if cs.Available {
		return ctx.Next(update)
	}

	slog.Debug("bot handler", "upd", update.UpdateID, "command", "cancel")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	_, _, cmdArgs := tu.ParseCommand(update.Message.Text)
	if len(cmdArgs) > 0 {
		_, err := bot.SendMessage(ctx, tu.Message(chatID, "Invalid amount of arguments in the command.\nThis is cancelation command that does not require any arguments.\nType:\n /cancel"))
		return err
	}
	if err := dbconvo.Free(update.Message.From.ID); err != nil {
		return err
	}
	_, err := bot.SendMessage(ctx, tu.Message(chatID, "Conversation was canceled, bot returned to regular working flow."))
	return err
}

func HandleConversation(ctx *th.Context, update telego.Update) error {
	cs := (ctx.Value("ConversationStatus")).(dbconvo.Status)
	if cs.Available {
		return ctx.Next(update)
	}

	if _, exists := CommandsList[cs.CommandName]; !exists {
		return fmt.Errorf("%s command does not exists in the commands list, but got in conversation status", cs.CommandName)
	} else if CommandsList[cs.CommandName].Conversation == nil {
		return fmt.Errorf("%s command conversation handler is nil, but got in conversation status", cs.CommandName)
	}
	return (*CommandsList[cs.CommandName].Conversation)(ctx, update)
}

func ConnectConversation(bh *th.BotHandler) {
	bh.Handle(HandleStatus, th.AnyMessage())
	bh.Handle(HandleCancel, th.CommandEqual("cancel"))
	bh.Handle(HandleConversation, th.AnyMessage())
}
