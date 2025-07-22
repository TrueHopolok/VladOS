package bot

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// Stores info about the user's engagement with the commands.
type ConversationStatus struct {
	// If user is free from any conversation.
	Free bool

	// Name of command for whom conversation is.
	Owner string
}

// TODO
func HandleConversation(ctx *th.Context, update telego.Update) error {
	return ctx.Next(update) // TODO: remove this for code to be executed
	// cs := ConversationStatus{} // TODO: GET FROM DB
	// if cs.Free {
	// 	return ctx.Next(update)
	// } else if _, exists := CommandsList[cs.Owner]; !exists {
	// 	return fmt.Errorf("%s command does not exists in the commands list, but got in conversation status", cs.Owner)
	// } else if CommandsList[cs.Owner].Conversation == nil {
	// 	return fmt.Errorf("%s command conversation handler is nil, but got in conversation status", cs.Owner)
	// }
	// return (*CommandsList[cs.Owner].Conversation)(ctx, update)
}

// TODO
//
//	in DB set [ConversationStatus.Free] to true
//	return nil gurantee
func HandleCancel(ctx *th.Context, update telego.Update) error {
	return nil
}

func ConnectConversation(bh *th.BotHandler) {
	bh.Handle(HandleCancel, th.CommandEqual("cancel"))
	bh.Handle(HandleConversation, th.AnyMessage())
}
