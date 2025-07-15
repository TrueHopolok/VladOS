package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TrueHopolok/VladOS/modules/bot"
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
	// opening
	slog.Info("db init", "status", "START")
	if err := db.Init(); err != nil {
		slog.Error("db init", "status", "FAILED", "error", err)
		return
	}
	// closing
	defer func() {
		slog.Info("db close", "status", "START")
		if err := db.Conn.Close(); err != nil {
			slog.Error("db close", "status", "FAILED", "error", err)
		} else {
			slog.Info("db close", "status", "SUCCESS")
		}
	}()
	// success
	slog.Info("db init", "status", "SUCCESS")

	//* DB migrate
	slog.Info("db migrate", "status", "START")
	if err := db.Migrate(); err != nil {
		slog.Error("db migrate", "status", "FAILED", "error", err)
		return
	}
	slog.Info("db migrate", "status", "SUCCESS")

	//* HTTP initialization
	// opening
	slog.Info("http init", "status", "START")
	httpErrorChan := make(chan error)
	if err := web.Start(httpErrorChan); err != nil {
		slog.Error("http init", "status", "FAILED", "error", err)
		return
	}
	// closing
	defer func() {
		slog.Info("http close", "status", "START")
		if err := web.Stop(); err != nil {
			slog.Error("http close", "status", "FAILED", "error", err)
		} else {
			slog.Info("http close", "status", "SUCCESS")
		}
	}()
	// success
	slog.Info("http init", "status", "SUCCESS")

	//* Bot initialization
	// opening
	slog.Info("bot init", "status", "SUCCESS")
	botErrorChan := make(chan error)
	if err := bot.Start(botErrorChan); err != nil {
		slog.Error("bot init", "status", "FAILED", "error", err)
		return
	}
	// closing
	defer func() {
		slog.Info("bot close", "status", "START")
		if err := bot.Stop(); err != nil {
			slog.Error("bot close", "status", "FAILED", "err", err)
		}
		slog.Info("bot close", "status", "SUCCESS")
	}()
	// success
	slog.Info("bot init", "status", "SUCCESS")

	//* Interupt initialization
	slog.Info("interupt init", "status", "START")
	sigChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}()
	slog.Info("interupt init", "status", "FINISH")

	//* Multithreads listening
	select {
	case err := <-httpErrorChan:
		if errors.Is(err, http.ErrServerClosed) {
			slog.Warn("http execute", "status", "STOPPED", "msg", "Server was closed without errors")
		} else {
			slog.Error("http execute", "status", "FAILED", "error", err)
		}
	case err := <-botErrorChan:
		slog.Error("bot execute", "status", "FAILED", "error", err)
	case <-sigChan:
		slog.Warn("interupt execute", "status", "CAUGHT")
	}
}
