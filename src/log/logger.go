package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var logger *zap.Logger

func init() {
	//w := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   "foo.log",
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 3,
	//	MaxAge:     28, // days
	//})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.DebugLevel,
	)
	logger = zap.New(core)

}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)

}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

type Log interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

type TLog struct {
	prefix string // 日志前缀
}

func NewTLog(prefix string) *TLog {

	return &TLog{prefix:prefix}
}

func (t *TLog) Info(msg string, fields ...zap.Field)  {
	Info(fmt.Sprintf("【%s】%s",t.prefix,msg),fields...)
}

func (t *TLog) Debug(msg string, fields ...zap.Field)  {
	Debug(fmt.Sprintf("【%s】%s",t.prefix,msg),fields...)
}
func (t *TLog) Error(msg string, fields ...zap.Field) {
	Error(fmt.Sprintf("【%s】%s",t.prefix,msg),fields...)
}
func (t *TLog) Warn(msg string, fields ...zap.Field)  {
	Warn(fmt.Sprintf("【%s】%s",t.prefix,msg),fields...)
}