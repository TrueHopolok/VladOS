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
	slog.Info("db init", "status", "start")
	if err := db.Init(); err != nil {
		slog.Error("db init", "status", "failed", "error", err)
	}
	slog.Info("db init", "status", "success")
}
