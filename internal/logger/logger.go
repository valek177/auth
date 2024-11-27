package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init initializes new logger
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// Debug is used for debug logging
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info is used for info logging
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn is used for warn logging
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error is used for error logging
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// ErrorWithMsg is used for error logging with error param
func ErrorWithMsg(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	globalLogger.Error(msg, fields...)
}

// Fatal is used for fatal logging
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// FatalWithMsg is used for fatal logging with error param
func FatalWithMsg(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	globalLogger.Fatal(msg, fields...)
}

// WithOptions applies options
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
