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
    Rotate int
    MaxSize int
    MaxBackup int
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
    if defaultLog != nil {
        defaultLog.Warn("can not install logger twice !!!")
        return defaultLog
    }

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
    defaultLog = &l
    return &l
}

func (l *LevelLogger) Log(level int, depth int, v ...interface{} ) {
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
        if len(v) == 1 {
            message = fmt.Sprintf("%+v", v[0])
        } else {
            message = fmt.Sprintf(v[0].(string), v[1:]...)
        }
        if l.Filename == "stdout" {
            // XXX debug only, slow with 4 lock
            l.mu.Lock()
            l.Logger.SetPrefix(color)
            l.Logger.Output(depth, fmt.Sprint(tag, " ", l.Prefix, message, colorReset))
            l.Logger.SetPrefix("")
            l.mu.Unlock()
        } else {
            l.Logger.Output(depth, fmt.Sprint(tag, " ", l.Prefix, message))
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
    l.Log(LevelDebug, 3, v...)
}

func (l *LevelLogger) Info(v ...interface{}) {
    l.Log(LevelInfo, 3, v...)
}

func (l *LevelLogger) Warn(v ...interface{}) {
    l.Log(LevelWarn, 3, v...)
}

func (l *LevelLogger) Error(v ...interface{}) {
    l.Log(LevelError, 3, v...)
}

func (l *LevelLogger) Fatal(v ...interface{}) {
    l.Log(LevelFatal, 3, v...)
    os.Exit(1)
}

func Debug(v ...interface{}) {
    defaultLog.Log(LevelDebug, 3, v...)
}

func Info(v ...interface{}) {
    defaultLog.Log(LevelInfo, 3, v...)
}

func Warn(v ...interface{}) {
    defaultLog.Log(LevelWarn, 3, v...)
}

func Error(v ...interface{}) {
    defaultLog.Log(LevelError, 3, v...)
}

func Fatal(v ...interface{}) {
    defaultLog.Log(LevelFatal, 3, v...)
    os.Exit(1)
}

func Debugd(depth int, v ...interface{}) {
    defaultLog.Log(LevelDebug, 3+depth, v...)
}

func Infod(depth int, v ...interface{}) {
    defaultLog.Log(LevelInfo, 3+depth, v...)
}

func Warnd(depth int, v ...interface{}) {
    defaultLog.Log(LevelWarn, 3+depth, v...)
}

func Errord(depth int, v ...interface{}) {
    defaultLog.Log(LevelError, 3+depth, v...)
}

func Fatald(depth int, v ...interface{}) {
    defaultLog.Log(LevelFatal, 3+depth, v...)
    os.Exit(1)
}
