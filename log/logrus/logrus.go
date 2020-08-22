package logrus

import (
	"bytes"
	"fmt"
	"github.com/dbunion/com/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

// Log is log adapter.
type Log struct {
	config    log.Config
	logger    *logrus.Logger
	outWriter io.Writer
}

// NewLogrus create new logrus log with default collection name.
func NewLogrus() log.Logger {
	return &Log{}
}

// Infof - info format log
func (l *Log) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Info - info format log
func (l *Log) Info(v ...interface{}) {
	l.logger.Info(v...)
}

// Debugf - debug format log
func (l *Log) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

// Debug - debug log
func (l *Log) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

// Warnf - warn format log
func (l *Log) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

// Warn - warn log
func (l *Log) Warn(v ...interface{}) {
	l.logger.Warn(v...)
}

// Warningf - Warning format log
func (l *Log) Warningf(format string, v ...interface{}) {
	l.logger.Warningf(format, v...)
}

// Warning - Warning log
func (l *Log) Warning(v ...interface{}) {
	l.logger.Warning(v...)
}

// Errorf - error format log
func (l *Log) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

// Error - error log
func (l *Log) Error(v ...interface{}) {
	l.logger.Error(v...)
}

// Fatalf - fatal format log
func (l *Log) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

// Fatal - fatal format log
func (l *Log) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

// Fatalln - fatal log
func (l *Log) Fatalln(v ...interface{}) {
	l.logger.Fatalln(v...)
}

// Printf - print format log
func (l *Log) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Print - print log
func (l *Log) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Println - print log
func (l *Log) Println(v ...interface{}) {
	l.logger.Println(v...)
}

// Panic - panic
func (l *Log) Panic(v ...interface{}) {
	l.logger.Panic(v...)
}

// Panicf - panic format value
func (l *Log) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

// Panicln - panic
func (l *Log) Panicln(v ...interface{}) {
	l.logger.Panicln(v...)
}

// Close connection
func (l *Log) Close() error {
	return nil
}

// StartAndGC start log adapter.
func (l *Log) StartAndGC(config log.Config) error {
	config.CheckWithDefault()

	l.config = config
	l.logger = logrus.New()

	l.logger.SetReportCaller(true)
	l.logger.SetLevel(l.getLogLevel())

	opts := []rotatelogs.Option{
		rotatelogs.WithLinkName(config.FilePath),
		rotatelogs.WithRotationTime(config.RotationTime),
	}

	if config.RotationMaxAge > 0 {
		opts = append(opts, rotatelogs.WithMaxAge(config.RotationMaxAge))
	}

	if config.RotationCount > 0 {
		opts = append(opts, rotatelogs.WithRotationCount(config.RotationCount))
	}

	writer, err := rotatelogs.New(l.config.FilePath+".%Y%m%d%H%M", opts...)
	if err != nil {
		return err
	}

	l.outWriter = writer
	l.logger.SetOutput(writer)
	l.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05.99",
		ForceColors:               config.HighLighting,
		EnvironmentOverrideColors: config.HighLighting,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			pc, file, line, _ := runtime.Caller(config.CallerSkip)
			name := runtime.FuncForPC(pc).Name()
			if i := bytes.LastIndexAny([]byte(name), "."); i != -1 {
				name = name[i+1:]
			}

			if i := bytes.LastIndexAny([]byte(file), "/"); i != -1 {
				file = file[i+1:]
			}

			return name, fmt.Sprintf("%s:%d", file, line)
		},
	})
	if config.JSONFormatter {
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.99",
			DisableTimestamp: false,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				pc, file, line, _ := runtime.Caller(9)
				name := runtime.FuncForPC(pc).Name()
				if i := bytes.LastIndexAny([]byte(name), "."); i != -1 {
					name = name[i+1:]
				}

				if i := bytes.LastIndexAny([]byte(file), "/"); i != -1 {
					file = file[i+1:]
				}

				return name, fmt.Sprintf("%s:%d", file, line)
			},
			PrettyPrint: false,
		})
	}
	return nil
}

func (l *Log) getLogLevel() logrus.Level {
	switch l.config.Level {
	case log.LevelInfo:
		return logrus.InfoLevel
	case log.LevelDebug:
		return logrus.DebugLevel
	case log.LevelWarning:
		return logrus.WarnLevel
	case log.LevelError:
		return logrus.ErrorLevel
	case log.LevelFatal:
		return logrus.FatalLevel
	}
	return logrus.InfoLevel
}

func init() {
	log.Register(log.TypeLogrus, NewLogrus)
}
