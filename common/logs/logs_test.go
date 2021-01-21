package logs

import (
	"errors"
	"testing"
)

func TestLogger_Init(t *testing.T) {
	var Silent bool
	var Level Level
	var Path string
	var FileStatus bool
	l := Logger{
		Silent:     &Silent,
		Level:      &Level,
		Path:       &Path,
		FileStatus: &FileStatus,
	}
	l.Init()
	l.Debug("Debug", "test", errors.New("l.Debug"))
	l.Debugf("%s %s %s", "Debug", "test", errors.New("l.Debugf"))
	//Level = ERROR
	l.Info("info", "test", errors.New("l.Info"))
	l.Infof("%s %s %s", "info", "test", errors.New("l.Infof"))
	l.Warmming("Warmming", "test", errors.New("l.Warmming"))
	l.Warmmingf("%s %s %s", "Warmming", "test", errors.New("l.Warmmingf"))
	//l.Exit("Exit", "test", errors.New("l.Exit"))
	l.Exitf("%s %s %s", "Exit", "test", errors.New("l.Exitf"))
}
