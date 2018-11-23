
package logger

import (
    "fmt"
    "io"
    "log"
    "sync"
)

const (
    levelDebug = (iota + 1) * 10
    levelInfo
    levelWarn
    levelError
    levelFatal
)

const (
    tagDebug = "D"
    tagInfo  = "I"
    tagWarn  = "W"
    tagError = "E"
    tagFatal = "F"
)

const (
    colorDebug = ""
    colorInfo = "\33[36m"
    colorWarn = "\33[33m"
    colorError = "\33[35m"
    colorFatal = "\33[31m"
    colorReset = "\33[0m"
)

const (
    flags = log.Ldate | log.Lmicroseconds | log.Lshortfile
)

var (
    logLock sync.Mutex
)
