package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

func init() {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	appLogger = logger.Sugar()
	defer appLogger.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func Debug(args ...interface{}) {
	appLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	appLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	appLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	appLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	appLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	appLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	appLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	appLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	appLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	appLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	appLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	appLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	appLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	appLogger.Fatalf(template, args...)
}
