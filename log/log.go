// TODO log time rotate
// TODO log size rotate
// TODO multi logger
// TODO windows color
// TODO check tty

package log

import (
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
)

const (
    LevelDebug = (iota + 1) * 10
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal
)

const (
    tagDebug = "[D]"
    tagInfo  = "[I]"
    tagWarn  = "[W]"
    tagError = "[E]"
    tagFatal = "[F]"
    tagLog   = "[L]"
)

// ref: https://en.wikipedia.org/wiki/ANSI_escape_code
const (
    colorDebug = "\033[37m"
    colorInfo  = "\033[36m"
    colorWarn  = "\033[33m"
    colorError = "\033[31m"
    colorFatal = "\033[35m"
    colorReset = "\033[0m"
)

// TODO
const (
    RotateNo = iota
    RotateTimeDay
    RotateTimeHour
    RotateTimeMinute
    RotateTimeSecond
    RotateSizeKB
    RotateSizeMB
    RotateSizeGB
)

type LevelLogger struct {
    *log.Logger
    fp *os.File
    mu sync.Mutex
    Prefix string
    Filename  string
    Level int

    // TODO
    Rotate int
    MaxSize int
    MaxBackup int
}

var (
    DefaultLog *LevelLogger
)

func GetLogger() *LevelLogger {
    if DefaultLog == nil {
        fmt.Println("can not GetLogger before Install")
        os.Exit(1)
    }
    return DefaultLog
}

func Install(dest string) *LevelLogger {

    var base *log.Logger
    var fp *os.File

    if dest == "stdout" {
        fp = os.Stdout
        base = log.New(fp, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    } else {
        fp, err := os.OpenFile(dest, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            fmt.Printf("can not open logfile: %v\n", err)
        }
        base = log.New(fp, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    }

    l := LevelLogger{
        Logger: base,
        fp: fp,
        Prefix: "",
        Filename:   dest,
        Level:  LevelDebug,
    }

    // first logger as DefaultLog
    if DefaultLog == nil {
        DefaultLog = &l
    }

    return &l
}

func (l *LevelLogger) Log(level int, depth int, prefix string, v ...interface{} ) {
    if len(v) == 0 {
        return
    }
    if level >= l.Level {
        var tag,color,message string
        switch level {
            case LevelDebug:
                tag = tagDebug
                color = colorDebug
            case LevelInfo:
                tag = tagInfo
                color = colorInfo
            case LevelWarn:
                tag = tagWarn
                color = colorWarn
            case LevelError:
                tag = tagError
                color = colorError
            case LevelFatal:
                tag = tagFatal
                color = colorFatal
            default:
                tag = tagLog
                color = colorReset
        }
        if format, ok := v[0].(string); ok {
            message = fmt.Sprintf(format, v[1:]...)
        } else {
            format := strings.Repeat("%+v ", len(v))
            message = fmt.Sprintf(format, v...)
        }
        if l.Filename == "stdout" {
            // XXX debug only, slow with 4 lock
            l.mu.Lock()
            l.Logger.SetPrefix(color)
            l.Logger.Output(depth, fmt.Sprint(tag, " ", prefix, message, colorReset))
            l.Logger.SetPrefix("")
            l.mu.Unlock()
        } else {
            l.Logger.Output(depth, fmt.Sprint(tag, " ", prefix, message))
        }
    }
}

func (l *LevelLogger) SetPrefix(prefix string) {
    l.Prefix = prefix
}

func (l *LevelLogger) SetLevel(level int) {
    l.Level = level
}

func (l *LevelLogger) Debug(v ...interface{}) {
    l.Log(LevelDebug, 3, l.Prefix, v...)
}

func (l *LevelLogger) Info(v ...interface{}) {
    l.Log(LevelInfo, 3, l.Prefix, v...)
}

func (l *LevelLogger) Warn(v ...interface{}) {
    l.Log(LevelWarn, 3, l.Prefix, v...)
}

func (l *LevelLogger) Error(v ...interface{}) {
    l.Log(LevelError, 3, l.Prefix, v...)
}

func (l *LevelLogger) Fatal(v ...interface{}) {
    l.Log(LevelFatal, 3, l.Prefix, v...)
    os.Exit(1)
}

func Debug(v ...interface{}) {
    DefaultLog.Log(LevelDebug, 3, DefaultLog.Prefix, v...)
}

func Info(v ...interface{}) {
    DefaultLog.Log(LevelInfo, 3, DefaultLog.Prefix, v...)
}

func Warn(v ...interface{}) {
    DefaultLog.Log(LevelWarn, 3, DefaultLog.Prefix, v...)
}

func Error(v ...interface{}) {
    DefaultLog.Log(LevelError, 3, DefaultLog.Prefix, v...)
}

func Fatal(v ...interface{}) {
    DefaultLog.Log(LevelFatal, 3, DefaultLog.Prefix, v...)
    os.Exit(1)
}

func Debugd(depth int, v ...interface{}) {
    DefaultLog.Log(LevelDebug, 3+depth, DefaultLog.Prefix, v...)
}

func Infod(depth int, v ...interface{}) {
    DefaultLog.Log(LevelInfo, 3+depth, DefaultLog.Prefix, v...)
}

func Warnd(depth int, v ...interface{}) {
    DefaultLog.Log(LevelWarn, 3+depth, DefaultLog.Prefix, v...)
}

func Errord(depth int, v ...interface{}) {
    DefaultLog.Log(LevelError, 3+depth, DefaultLog.Prefix, v...)
}

func Fatald(depth int, v ...interface{}) {
    DefaultLog.Log(LevelFatal, 3+depth, DefaultLog.Prefix, v...)
    os.Exit(1)
}
