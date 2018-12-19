// TODO windows color
// TODO check tty

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
    ColorDebug = "\033[37m"
    ColorInfo  = "\033[36m"
    ColorWarn  = "\033[33m"
    ColorError = "\033[31m"
    ColorFatal = "\033[35m"
    ColorReset = "\033[0m"
)


type LevelLogger struct {
    *log.Logger
    prefix string
    dest  string
    level int
    isColor bool
}

var (
    defaultLog *LevelLogger
)

func GetLogger() *LevelLogger {
    if defaultLog == nil {
        fmt.Println("can not GetLogger before Install")
        os.Exit(1)
    }
    return defaultLog
}

func Install(dest string) *LevelLogger {
    var isColor bool
    var base *log.Logger

    if dest == "stdout" {
        isColor = true
        out := os.Stdout
        base = log.New(out, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    } else {
        isColor = false
        out, err := os.OpenFile(dest, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            fmt.Printf("can not open logfile: %v\n", err)
        }
        base = log.New(out, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    }

    l := LevelLogger{
        Logger: base,
        prefix: "",
        dest:   dest,
        level:  LevelDebug,
        isColor: isColor,
    }
    defaultLog = &l
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
            l.Logger.SetPrefix(color)
            l.Logger.Output(3, fmt.Sprint(tag, " ", l.prefix, message, ColorReset))
            l.Logger.SetPrefix("")
        } else {
            l.Logger.Output(3, fmt.Sprint(tag, " ", l.prefix, message))
        }
    }
}

func (l *LevelLogger) SetPrefix(prefix string) {
    l.prefix = prefix
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
    os.Exit(1)
}
