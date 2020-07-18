package log

import (
	"testing"
)

func TestLogOut(t *testing.T) {
	Debug("Debug")
	Error("Error")
	SetLevel(-8)
	Info("Info")
	Warn("Warn")
	Fatal("Fatal")

	defaultLogger := NewLogger("stderr", -2, 0, true)
	defaultLogger.Debug("Debug")
	defaultLogger.Error("Error")
	defaultLogger.SetLevel(-8)
	defaultLogger.Info("Info")
	defaultLogger.Warn("Warn")
	defaultLogger.Fatal("Fatal")
	defaultLogger.rotate()
	defaultLogger.Close()

	defaultLogger = NewLogger("err", 1, 0, true)
	defaultLogger.Debug("Debug")
	defaultLogger.Error("Error")
	defaultLogger.Info("Info")
	defaultLogger.Warn("Warn")
	defaultLogger.Trace("Trace")
	defaultLogger.rotate()
	defaultLogger.Close()

	OpenDefaultLog("./err.20200717-161712", LevelDebug, 0, true)
}
