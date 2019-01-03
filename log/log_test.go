package log

import "errors"
import "testing"

func TestInstall(t *testing.T) {
    Install("stdout")
    log := GetLogger()
    log.Debug("test")
    log.Debug("test %s", "format")

    Debug(errors.New("test error"))
    Debug("test")
    Info("test %s", "format")
    Warn("test %s", "format")
    Error("test %s", "format")
}
func TestLevel(t *testing.T) {
    Install("stdout")
    log := GetLogger()
    for i := 0; i<10 ; i++ {
        log.Debug("中文 debug %d", i)
        log.Info("😀 info %d", i)
        log.Warn("warn %d", i)
        log.Error("error %d", i)
        log.Printf("print %d", i)
    }

    log.SetLevel(LevelWarn)
    for i := 0; i<10 ; i++ {
        log.Debug("debug %d", i)
        log.Info("info %d", i)
        log.Warn("warn %d", i)
        log.Error("error %d", i)
    }
    // log.Fatal("fatal")
}
func TestPrefix(t *testing.T) {
    log := Install("stdout")
    log.Debug("test")
    log.SetPrefix("prefix:")
    log.Debug("test")
    log.Debug("test")
    log.SetPrefix("")
    log.Debug("test")
}
