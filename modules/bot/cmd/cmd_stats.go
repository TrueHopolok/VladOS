package cmd

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

// TODO: make whole stats and update on for each new machine (potentialy make it dynamic)
var CommandStats Command = Command{
	InfoBrief: "output stats for game",
	InfoFull: `
 /stats <game_name>
Output all stats for ceratin game.
With a your placement and top placement in the leaderboard.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		return nil
	},
	conversation: nil,
}
