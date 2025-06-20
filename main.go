package main

import (
	"log/slog"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/mlog"
)

func main() {
	//* Logger initialization
	mlog.Init()
	defer func() {
		if x := recover(); x != nil {
			slog.Error("main caught panic, stopping execution", "panic", x)
		}
	}()
	slog.Info("program started")

	//* DB initialization
	slog.Info("db init", "STATUS", "START")
	if err := db.Init(); err != nil {
		slog.Error("db init", "STATUS", "FAILED", "error", err)
		return
	}
	slog.Info("db init", "STATUS", "SUCCESS")

	//* DB migrate
	slog.Info("db migrate", "STATUS", "START")
	if err := db.Migrate(); err != nil {
		slog.Error("db migrate", "STATUS", "FAILED", "error", err)
		return
	}
	slog.Info("db migrate", "STATUS", "SUCCESS")
}
