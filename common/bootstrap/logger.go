package bootstrap

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
)

// Log 等级
const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

// Logger Log 包装
type Logger struct {
	*logrus.Logger
}

// Log 全局 Log
var Log *Logger

// NewLogger 初始化 Logger
func NewLogger() (*Logger, *errors.Error) {
	log := &Logger{
		Logger: logrus.New(),
	}

	log.Logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	}

	if config.Config.Dev {
		log.Logger.SetLevel(LogLevel(DebugLevel))
	} else {
		log.Logger.SetLevel(LogLevel(InfoLevel))
	}

	log.Logger.Out = os.Stdout
	return log, nil
}

// Logf 根据等级输出
func (log *Logger) Logf(level logrus.Level, format string, argu ...interface{}) {
	type loggerFunc func(string, ...interface{})

	var loggerAction = map[logrus.Level]loggerFunc{
		PanicLevel: log.Panicf,
		FatalLevel: log.Panicf,
		ErrorLevel: log.Errorf,
		WarnLevel:  log.Warnf,
		InfoLevel:  log.Infof,
		DebugLevel: log.Debugf,
	}

	loggerAction[level](format, argu...)
}

// LogLevel 返回 Log 等级
func LogLevel(defaultLevel logrus.Level) logrus.Level {
	levels := map[string]logrus.Level{
		"panic": PanicLevel,
		"fatal": FatalLevel,
		"error": ErrorLevel,
		"warn":  WarnLevel,
		"info":  InfoLevel,
		"debug": DebugLevel,
	}

	if level, ok := levels[config.Config.LogLevel]; ok {
		return level
	}

	return defaultLevel
}
