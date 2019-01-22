package log

import "testing"

func TestContext(t *testing.T) {
    Install("stdout")

    ctx := NewLoggerContext(nil, "prefix: ")
    ctx.Debug("debug")
    ctx.Info("info")
    ctx.Warn("warn")
    ctx.Error("warn")

    ctx.SetValue("key", "test value")
    v := ctx.Value("key").(string)
    ctx.Debug(v)
}
