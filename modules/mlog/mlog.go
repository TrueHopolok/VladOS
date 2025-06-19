// mlog package shorthanded for MyLogger.
//
// Contains all necessary functional to begin work with [log/slog] package.
// Uses [github.com/TrueHopolok/VladOS/modules/cfg] package to get all necessary info for initialization.
package mlog

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/cfg"

	"gopkg.in/natefinch/lumberjack.v2"
)

const LogFilePath string = "logs/"

// Initializate the slog package for performing purposes with given config parameters by calling [github.com/TrueHopolok/VladOS/cfg.Get].
// Uses [github.com/TrueHopolok/VladOS/cfg.Config.Verbose] to decide wheter to output the debug prints (in case of true value) or ignore them.
// Creates multiwriter to print in both [os.Stdout] and [github.com/TrueHopolok/VladOS/cfg.Config.LogFilePath].
// Creates writer that uses [gopkg.in/natefinch/lumberjack.v2] package to perform the log file rotations.
// Writer is set as default output handler of [log/slog] package, to be used outside this package.
func Init() {
	writer := &lumberjack.Logger{
		Filename: LogFilePath + cfg.Get().LogFileName,
		MaxSize:  cfg.Get().LogMaxSize,
	}

	logLevel := slog.LevelInfo
	if cfg.Get().Verbose {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(
		io.MultiWriter(writer, os.Stdout),
		&slog.HandlerOptions{
			Level: logLevel,
		})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Debug("logger was initialized regularly")
}

// Initializate the slog package for performing purposes with given config parameters by calling [github.com/TrueHopolok/VladOS/cfg.Get].
// Ignores [github.com/TrueHopolok/VladOS/cfg.Config.Verbose] and always prints all [log/slog.Level] including debug one.
// Does not print in [os.Stdout] and only prints in [github.com/TrueHopolok/VladOS/cfg.Config.LogFilePath].
// Creates writer that uses [gopkg.in/natefinch/lumberjack.v2] package to perform the log file rotations.
// Writer is set as default output handler of [log/slog] package, to be used outside this package.
func InitTesting(t *testing.T, pathToRoot string) {
	if !testing.Testing() {
		panic(fmt.Errorf("tried to initialize the logger in test mode while not in testing mode"))
	}

	writer := &lumberjack.Logger{
		Filename: pathToRoot + LogFilePath + cfg.TestGet(pathToRoot).LogFileName,
		MaxSize:  cfg.TestGet(pathToRoot).LogMaxSize,
	}

	handler := slog.NewTextHandler(
		writer,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Debug("logger was initialized for testing")
}
