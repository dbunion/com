package zssky

import (
	"github.com/dbunion/com/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	zslog "github.com/zssky/log"
)

// Log is log adapter.
type Log struct {
	config log.Config
}

// NewZsskyLog create new zssky log with default collection name.
func NewZsskyLog() log.Logger {
	return &Log{}
}

// Infof - info format log
func (l *Log) Infof(format string, v ...interface{}) {
	zslog.Infof(format, v...)
}

// Info - info format log
func (l *Log) Info(v ...interface{}) {
	zslog.Info(v...)
}

// Debugf - debug format log
func (l *Log) Debugf(format string, v ...interface{}) {
	zslog.Debugf(format, v...)
}

// Debug - debug log
func (l *Log) Debug(v ...interface{}) {
	zslog.Debug(v...)
}

// Warnf - warn format log
func (l *Log) Warnf(format string, v ...interface{}) {
	zslog.Warnf(format, v...)
}

// Warn - warn log
func (l *Log) Warn(v ...interface{}) {
	zslog.Warn(v...)
}

// Warningf - Warning format log
func (l *Log) Warningf(format string, v ...interface{}) {
	zslog.Warningf(format, v...)
}

// Warning - Warning log
func (l *Log) Warning(v ...interface{}) {
	zslog.Warning(v...)
}

// Errorf - error format log
func (l *Log) Errorf(format string, v ...interface{}) {
	zslog.Errorf(format, v...)
}

// Error - error log
func (l *Log) Error(v ...interface{}) {
	zslog.Error(v...)
}

// Fatalf - fatal format log
func (l *Log) Fatalf(format string, v ...interface{}) {
	zslog.Fatalf(format, v...)
}

// Fatal - fatal format log
func (l *Log) Fatal(v ...interface{}) {
	zslog.Fatal(v...)
}

// Fatalln - fatal log
func (l *Log) Fatalln(v ...interface{}) {
	zslog.Fatal(v...)
}

// Printf - print format log
func (l *Log) Printf(format string, v ...interface{}) {
	zslog.Infof(format, v...)
}

// Print - print format log
func (l *Log) Print(v ...interface{}) {
	zslog.Info(v...)
}

// Println - print value
func (l *Log) Println(v ...interface{}) {
	zslog.Info(v...)
}

// Panic - panic
func (l *Log) Panic(v ...interface{}) {
	zslog.Error(v...)
}

// Panicf - panic format value
func (l *Log) Panicf(format string, v ...interface{}) {
	zslog.Errorf(format, v...)
}

// Panicln - panic
func (l *Log) Panicln(v ...interface{}) {
	zslog.Error(v...)
}

// Close connection
func (l *Log) Close() error {
	return nil
}

// StartAndGC start log adapter.
func (l *Log) StartAndGC(config log.Config) error {
	config.CheckWithDefault()

	l.config = config

	zslog.SetLevelByString(string(config.Level))

	// basic setting
	zslog.SetHighlighting(config.HighLighting)

	// rotate opts
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

	zslog.SetOutput(writer)

	callerSkip := 5
	if config.CallerSkip != 0 {
		callerSkip = config.CallerSkip
	}
	zslog.SetCallerSkip(callerSkip)

	return nil
}

func init() {
	log.Register(log.TypeZsskyLog, NewZsskyLog)
}
