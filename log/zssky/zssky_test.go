package zssky

import (
	"github.com/dbunion/com/log"
	"testing"
	"time"
)

func TestZsskyInfoLog(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelInfo,
		FilePath:      "/tmp/zssky.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  true,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Infof("log test, date:%v", time.Now().Unix())
}

func TestZsskyDebugLog(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelDebug,
		FilePath:      "/tmp/zssky.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  true,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Debugf("log test, date:%v", time.Now().Unix())
}

func TestZsskyWarningLog(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelWarning,
		FilePath:      "/tmp/zssky.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  true,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Warningf("log test, date:%v", time.Now().Unix())
}

func TestZsskyErrorLog(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelError,
		FilePath:      "/tmp/zssky.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  true,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Errorf("log test, date:%v", time.Now().Unix())
}

func TestZsskyFatalLog(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelFatal,
		FilePath:      "/tmp/zssky.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  true,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Fatalf("log test, date:%v", time.Now().Unix())
}
