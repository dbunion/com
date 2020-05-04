package logrus

import (
	"github.com/dbunion/com/log"
	"testing"
	"time"
)

func TestLogrusInfo(t *testing.T) {
	logger, err := log.NewLogger(log.TypeLogrus, log.Config{
		Level:         log.LevelInfo,
		FilePath:      "/tmp/logrus.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  false,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Infof("logrus test, date:%v", time.Now().Unix())
}

func TestLogrusDebug(t *testing.T) {
	logger, err := log.NewLogger(log.TypeLogrus, log.Config{
		Level:         log.LevelDebug,
		FilePath:      "/tmp/logrus.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  false,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Debugf("logrus test, date:%v", time.Now().Unix())
}

func TestLogrusWarning(t *testing.T) {
	logger, err := log.NewLogger(log.TypeLogrus, log.Config{
		Level:         log.LevelWarning,
		FilePath:      "/tmp/logrus.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  false,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Warningf("logrus test, date:%v", time.Now().Unix())
}

func TestLogrusError(t *testing.T) {
	logger, err := log.NewLogger(log.TypeLogrus, log.Config{
		Level:         log.LevelError,
		FilePath:      "/tmp/logrus.log",
		HighLighting:  true,
		RotateByDay:   false,
		RotateByHour:  false,
		JSONFormatter: false,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	logger.Errorf("logrus test, date:%v", time.Now().Unix())
}
