package bot

import (
	"log/slog"

	spch "github.com/TrueHopolok/spellchecker/spellchecker"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleSpelling(ctx *th.Context, update telego.Update) error {
	slog.Debug("bot handler", "upd", update.UpdateID, "command", "spelling")
	bot := ctx.Bot()
	chatID := update.Message.Chat.ChatID()
	cmdName, _, _ := tu.ParseCommand(update.Message.Text)
	var msgText []tu.MessageEntityCollection
	msgText = append(msgText, tu.Entityf("'%s' is not a command.\nSee /help for whole list of the commands.\n\nThe most similar commands are:", cmdName))
	for _, potName := range CmdFindClosest(cmdName) {
		msgText = append(msgText, tu.Entityf("\n /%s", potName))
	}
	_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
	return err
}

func CmdFindClosest(cmdName string) []string {
	var potentialCommands []string
	minScore := 0
	for existingName := range CommandsList.Gambling {
		curScore := spch.FindScore(cmdName, existingName)
		if minScore == 0 || minScore > curScore {
			potentialCommands = potentialCommands[:0]
			minScore = curScore
		}
		if minScore == curScore {
			potentialCommands = append(potentialCommands, existingName)
		}
	}
	for existingName := range CommandsList.Others {
		curScore := spch.FindScore(cmdName, existingName)
		if minScore == 0 || minScore > curScore {
			potentialCommands = potentialCommands[:0]
			minScore = curScore
		}
		if minScore == curScore {
			potentialCommands = append(potentialCommands, existingName)
		}
	}
	extraCmdList := []string{"help", "start"}
	for _, potName := range extraCmdList {
		curScore := spch.FindScore(cmdName, potName)
		if minScore == 0 || minScore > curScore {
			potentialCommands = potentialCommands[:0]
			minScore = curScore
		}
		if minScore == curScore {
			potentialCommands = append(potentialCommands, potName)
		}
	}
	return potentialCommands
}
