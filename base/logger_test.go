package logger

import "testing"

func TestInstall(t *testing.T) {
    log := Install("stdout")
    log.Debug("test")
    log.Debug("test %s", "format")
}
func TestLevel(t *testing.T) {
    log := Install("stdout")
    for i := 0; i<10 ; i++ {
        log.Debug("中文 text %d", i)
        log.Info("😀 text %d", i)
        log.Warn("text %d", i)
        log.Error("text %d", i)
    }

    log.SetLevel(LevelWarn)
    for i := 0; i<10 ; i++ {
        log.Debug("text %d", i)
        log.Info("text %d", i)
        log.Warn("text %d", i)
        log.Error("text %d", i)
    }
}
