// [vlados/modules/mlog] package shorthanded for MyLogger.
//
// Contains all necessary functional to begin work with [log/slog] package.
// Uses [vlados/modules/cfg] package to get all necessary info for initialization.
package mlog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"vlados/modules/cfg"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Initializate the slog package for performing purposes with given config parameters by calling [vlados/modules/cfg.Get].
// Uses [vlados/modules/cfg.Config.Verbose] to decide wheter to output the debug prints (in case of true value) or ignore them.
// Creates multiwriter to print in both [os.Stdout] and [vlados/modules/cfg.Config.LogFilePath].
// Creates writer that uses [gopkg.in/natefinch/lumberjack.v2] package to perform the log file rotations.
// Writer is set as default output handler of [log/slog] package, to be used outside this package.
func RegularInit() {
	writer := &lumberjack.Logger{
		Filename: cfg.Get().LogFilePath,
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
	slog.Info("logger is inited via RegularInit")
}

// Initializate the slog package for testing purposes with given config parameters by calling [vlados/modules/cfg.Get].
// Ignores [vlados/modules/cfg.Config.Verbose] flag and outputs all level logs.
// Does not print in [os.Stdout] and only prints into [vlados/modules/cfg.Config.LogFilePath].
// Creates writer that uses [gopkg.in/natefinch/lumberjack.v2] package to perform the log file rotations.
// Writer is set as default output handler of [log/slog] package, to be used outside this package.
func TestingInit(t *testing.T) {
	if !testing.Testing() {
		panic(fmt.Errorf("tried to initialize the logger in test mode while not in testing mode"))
	}
	writer := &lumberjack.Logger{
		Filename: cfg.Get().LogFilePath,
		MaxSize:  cfg.Get().LogMaxSize,
	}

	handler := slog.NewTextHandler(
		writer,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("logger is inited via TestingInit")
}
