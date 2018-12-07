package logger

import (
    "fmt"
    "log"
    "os"
)

const (
    LevelDebug = (iota + 1) * 10
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal
)

const (
    TagDebug = "[D]"
    TagInfo  = "[I]"
    TagWarn  = "[W]"
    TagError = "[E]"
    TagFatal = "[F]"
    TagLog   = "[L]"
)

// ref: https://en.wikipedia.org/wiki/ANSI_escape_code
const (
    ColorDebug = ""
    ColorInfo  = "\033[36m"
    ColorWarn  = "\033[33m"
    ColorError = "\033[35m"
    ColorFatal = "\033[31m"
    ColorReset = "\033[0m"
)

type LevelLogger struct {
    *log.Logger
    dest  string
    level int
    isColor bool
}

func Install(dest string) *LevelLogger {
    var isColor bool

    if dest == "stdout" {
        isColor = true
    } else {
        isColor = false
    }

    l := LevelLogger{
        Logger: log.New(os.Stdout, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
        dest:   dest,
        level:  LevelDebug,
        isColor: isColor,
    }
    return &l
}

func (l *LevelLogger) Log(level int, format string, v ...interface{} ) {
    if level >= l.level {
        var tag,color,message string
        switch level {
            case LevelDebug:
                tag = TagDebug
                color = ColorDebug
            case LevelInfo:
                tag = TagInfo
                color = ColorInfo
            case LevelWarn:
                tag = TagWarn
                color = ColorWarn
            case LevelError:
                tag = TagError
                color = ColorError
            case LevelFatal:
                tag = TagFatal
                color = ColorFatal
            default:
                tag = TagLog
                color = ColorReset
        }
        if len(v) == 0 {
            message = format
        } else {
            message = fmt.Sprintf(format, v...)
        }
        if l.isColor{
            l.Logger.Output(3, fmt.Sprintln(color, tag, message, ColorReset))
        } else {
            l.Logger.Output(3, fmt.Sprintln(tag, message))
        }
    }
}

func (l *LevelLogger) SetLevel(level int) {
    l.level = level
}

func (l *LevelLogger) Debug(format string, v ...interface{}) {
    l.Log(LevelDebug, format, v...)
}

func (l *LevelLogger) Info(format string, v ...interface{}) {
    l.Log(LevelInfo, format, v...)
}

func (l *LevelLogger) Warn(format string, v ...interface{}) {
    l.Log(LevelWarn, format, v...)
}

func (l *LevelLogger) Error(format string, v ...interface{}) {
    l.Log(LevelError, format, v...)
}

func (l *LevelLogger) Fatal(format string, v ...interface{}) {
    l.Log(LevelFatal, format, v...)
}
