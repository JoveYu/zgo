
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

// ref: https://en.wikipedia.org/wiki/ANSI_escape_code
const (
    colorDebug = ""
    colorInfo  = "\33[36m"
    colorWarn  = "\33[33m"
    colorError = "\33[35m"
    colorFatal = "\33[31m"
    colorReset = "\33[0m"
)

const (
    Ldate = 1 << iota
    Ltime
    Lmicroseconds
    Llongfile
    Lshortfile
    LUTC
    Lpid
    Ltag
    LstdFlags = Ldate | Lmicroseconds | Lpid | Lshortfile | Ltag
)

type Logger struct {
    lock sync.Mutex
    level int
    dest string
    prefix string
    flag int
    out io.Writer
    buf []byte
}

