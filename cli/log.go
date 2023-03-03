package cli

import (
	"os"

	"github.com/benchkram/errz"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// This is the logger we can pass into services
var log logr.Logger

// logInit initializes the logger
func logInit(verbosity int8, pretty bool) {

	// assure info logs start at 1
	if verbosity != 0 {
		verbosity = verbosity - 1
	}

	// First, define our level-handling logic.
	stdError := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	stdOutput := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(-verbosity) && lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	if pretty {
		consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, stdError),
		zapcore.NewCore(consoleEncoder, consoleDebugging, stdOutput),
	)

	// multiline stacktrace option removed
	options := []zap.Option{} // zap.AddStacktrace(stdError),

	if pretty {
		options = append(options, zap.Development(), zap.AddStacktrace(stdError))
	}

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core, options...)

	log = zapr.NewLogger(logger)

	// Set the logger also for our error package
	errz.WithLogger(log.WithCallDepth(2))
}
