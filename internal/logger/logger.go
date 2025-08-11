package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = "development"
	}

	var encoder zapcore.Encoder
	var level zapcore.LevelEnabler

	if env == "development" {
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder // Adds color to level
		config.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		config.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewConsoleEncoder(config)
		level = zapcore.DebugLevel
	} else {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewJSONEncoder(config)
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
