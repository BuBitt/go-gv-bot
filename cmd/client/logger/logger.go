package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = customTimeEncoderFile
	fileEncoder := zapcore.NewJSONEncoder(config)

	config.EncodeTime = customTimeEncoder
	config.EncodeLevel = customLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile("logs/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// customTimeEncoderFile changes the timestamp format
func customTimeEncoderFile(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// customTimeEncoder changes the timestamp format to a custom layout with bold grey text
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("\x1b[1;30m%s\x1b[0m", t.Format("2006-01-02 15:04:05")))
}

// customLevelEncoder changes the level color
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	colorCode := levelColor(level)
	enc.AppendString(fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", colorCode, level.CapitalString()))
}

// levelColor returns the ANSI color code for the log level
func levelColor(level zapcore.Level) int {
	switch level {
	case zapcore.DebugLevel:
		return 36 // Cyan
	case zapcore.InfoLevel:
		return 32 // Green
	case zapcore.WarnLevel:
		return 33 // Yellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return 31 // Red
	default:
		return 0 // Default color
	}
}
