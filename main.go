package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	//* HTTP initialization
	slog.Info("http init", "status", "START")
	server := &http.Server{
		Addr:    ":8080",
		Handler: web.ConnectAll(),
	}
	httpErrorChan := make(chan error)
	go func() {
		httpErrorChan <- server.ListenAndServe()
	}()
	select {
	case err := <-httpErrorChan:
		slog.Error("http init", "status", "FAILED", "error", err)
		return
	default:
		slog.Info("http init", "status", "SUCCESS")
	}

	//* Bot initialization
	// TODO

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
	case <-sigChan:
		slog.Warn("interupt execute", "status", "CAUGHT")
	}

	//* Stop program's execution
	// bot
	// TODO

	// http
	slog.Info("http close", "status", "START")
	if err := server.Close(); err != nil {
		slog.Error("http close", "status", "FAILED", "error", err)
	} else {
		slog.Info("http close", "status", "SUCCESS")
	}

	// db
	slog.Info("db close", "status", "START")
	if err := db.Conn.Close(); err != nil {
		slog.Error("db close", "status", "FAILED", "error", err)
	} else {
		slog.Info("db close", "status", "SUCCESS")
	}

	// exiting the program
	os.Exit(0)
}
