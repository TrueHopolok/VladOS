package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/mlog"
	"github.com/TrueHopolok/VladOS/modules/web"
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
	slog.Info("db init", "status", "START")
	if err := db.Init(); err != nil {
		slog.Error("db init", "status", "FAILED", "error", err)
		return
	}
	slog.Info("db init", "status", "SUCCESS")

	//* DB migrate
	slog.Info("db migrate", "status", "START")
	if err := db.Migrate(); err != nil {
		slog.Error("db migrate", "status", "FAILED", "error", err)
		return
	}
	slog.Info("db migrate", "status", "SUCCESS")

	//* HTTP connect
	slog.Info("http", "status", "START")
	httpErrorChan := make(chan error)
	go func() {
		httpErrorChan <- http.ListenAndServe(":8080", web.ConnectAll())
	}()
	slog.Info("http", "status", "SUCCESS")

	//* Bot connect
	// TODO

	//* Console connect
	// TODO

	//* Multithreads listening
	select {
	case err := <-httpErrorChan:
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("http", "status", "STOPPED", "msg", "Server was closed without errors")
		} else {
			slog.Error("http", "status", "FAILED", "error", err)
		}
	}
}
