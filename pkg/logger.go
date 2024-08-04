package pkg

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger(filepath, stage string) (*zap.Logger, error) {
	consoleCore, err := createConsoleCore(stage)
	if err != nil {
		return nil, err
	}

	fileCore, err := createFileCore(filepath, stage)
	if err != nil {
		return nil, err
	}

	cores := zapcore.NewTee(
		consoleCore, fileCore,
	)

	return zap.New(cores), nil
}

func createConsoleCore(stage string) (zapcore.Core, error) {
	stdout := zapcore.AddSync(os.Stdout)
	cfg := zap.NewProductionEncoderConfig()
	setLogKeys(&cfg)

	enc := zapcore.NewConsoleEncoder(cfg)
	logLevel := chooseLogLevel(stage)

	return zapcore.NewCore(enc, stdout, logLevel), nil
}

func setLogKeys(cfg *zapcore.EncoderConfig) {
	cfg.TimeKey = "ts"
	cfg.CallerKey = "caller"
	cfg.LevelKey = "level"
	cfg.MessageKey = "msg"

	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.EpochMillisTimeEncoder
}

func chooseLogLevel(stage string) zapcore.Level {
	if stage == "prod" {
		return zap.InfoLevel
	}

	return zap.DebugLevel
}

func createFileCore(filepath, stage string) (zapcore.Core, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	fileOut := zapcore.AddSync(f)

	cfg := zap.NewProductionEncoderConfig()
	setLogKeys(&cfg)

	enc := zapcore.NewJSONEncoder(cfg)
	logLevel := chooseLogLevel(stage)

	return zapcore.NewCore(enc, fileOut, logLevel), nil
}
