package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitGlobalLogger(logFilePath, levelFlag string) error {
	loggingOutput := []string{"stdout"}
	errorOutput := []string{"stderr"}

	if len(logFilePath) > 0 {
		loggingOutput = append(loggingOutput, logFilePath)
		errorOutput = append(errorOutput, logFilePath)
	}

	logLevel := zapcore.DebugLevel

	switch levelFlag {
	case "info":
		logLevel = zapcore.InfoLevel
		break
	case "error":
		logLevel = zapcore.ErrorLevel
		break
	case "fatal":
		logLevel = zapcore.FatalLevel
		break
	case "warning":
		logLevel = zapcore.WarnLevel
		break
	case "debug":
		logLevel = zapcore.DebugLevel
		break
	default:
		logLevel = zapcore.InfoLevel
		break
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      loggingOutput,
		ErrorOutputPaths: errorOutput,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()

	if err != nil {
		log.Fatalf("Unable to initalise zap logger: %v ", err)
	}

	logger.Sugar().Debugf("Initialised logger with %s level and logging to %s", levelFlag, loggingOutput)

	zap.ReplaceGlobals(logger)

	return nil
}
