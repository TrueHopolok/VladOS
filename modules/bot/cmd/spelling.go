package cmd

import (
	"log/slog"

	spch "github.com/TrueHopolok/spellchecker/spellchecker"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func handleSpelling(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "spelling")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	cmdName, _, _ := tu.ParseCommand(update.Message.Text)
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entityf("'%s' is not a command.\nSee /help for whole list of the commands.\n\nThe most similar commands are:", cmdName))
	for _, potName := range utilClosestSpelling(cmdName) {
		msgText = append(msgText, tu.Entityf("\n /%s", potName))
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
	return err
}

// Given command name, finds closest existing commands by spelling.
func utilClosestSpelling(cmdName string) []string {
	var potentialCommands []string
	minScore := 0
	for category := range CommandsList {
		for existingName := range CommandsList[category] {
			curScore := spch.FindScore(cmdName, existingName)
			if minScore == 0 || minScore > curScore {
				potentialCommands = potentialCommands[:0]
				minScore = curScore
			}
			if minScore == curScore {
				potentialCommands = append(potentialCommands, existingName)
			}
		}
	}
	return potentialCommands
}
