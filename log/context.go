package log

import (
	"context"
)

type LoggerContext struct {
	context.Context
	Logger *LevelLogger

	// use this instead of Logger.Prefix
	Prefix string
}

func (c *LoggerContext) Debug(v ...interface{}) {
	c.Logger.Log(LevelDebug, 3, c.Prefix, v...)
}

func (c *LoggerContext) Info(v ...interface{}) {
	c.Logger.Log(LevelInfo, 3, c.Prefix, v...)
}

func (c *LoggerContext) Warn(v ...interface{}) {
	c.Logger.Log(LevelWarn, 3, c.Prefix, v...)
}

func (c *LoggerContext) Error(v ...interface{}) {
	c.Logger.Log(LevelError, 3, c.Prefix, v...)
}

func (c *LoggerContext) SetValue(k interface{}, v interface{}) {
	c.Context = context.WithValue(c, k, v)
}

func NewLoggerContext(logger *LevelLogger, prefix string) *LoggerContext {
	if logger == nil {
		logger = DefaultLog
	}
	return &LoggerContext{
		Logger: logger,
		Prefix: prefix,
	}
}
