
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
    tagDebug = "[D] "
    tagInfo = "[I] "
    tagWarn = "[W] "
    tagError = "[E] "
    tagFatal = "[F] "
)

const (
    flags = log.Ldate | log.Lmicroseconds | log.Lshortfile
)

var (
    logLock sync.Mutex
)
